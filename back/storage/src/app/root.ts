import express from "express";
import cors from "cors";
import logger from "morgan";
import cookieParser from "cookie-parser";
import { GetObjectCommand, S3Client } from "@aws-sdk/client-s3";
import { Readable } from "stream";

import upload from "../multerHandler";
import { uploadMain } from "./minio";
import { unlink } from "fs/promises";

require("dotenv").config();
const app = express();

// 環境変数の読み込み
const HOST_URL = process.env.HOST_URL || "";

const corsOptions = {
  origin: [HOST_URL],
  methods: "GET,HEAD,PUT,PATCH,POST,DELETE",
  optionsSuccessStatus: 200,
};

app.use(cors(corsOptions));
app.use(logger("dev")); // 開発用にログの表示
app.use(express.json()); // application/jsonを扱えるようにする
app.use(express.urlencoded({ extended: true }));
app.use(cookieParser());

const s3 = new S3Client({
  region: process.env.AWS_REGION || "ap-northeast-1",
  credentials: {
    accessKeyId: process.env.AWS_ACCESS_KEY_ID!,
    secretAccessKey: process.env.AWS_SECRET_ACCESS_KEY!,
  },
  endpoint: process.env.S3_ENDPOINT,
  forcePathStyle: true,
});

app.get("/ping", (req, res) => {
  res.json({ message: "connected from pictogram" });
});

app.post("/", upload.array("images", 7), async (req, res) => {
  if (!req.files || (req.files as Express.Multer.File[]).length === 0) {
    return res.status(400).send("No files uploaded");
  }

  const uid = req.body.uid;
  if (!uid) {
    for (const file of req.files as Express.Multer.File[]) {
      await unlink(file.path);
    }
    return res.status(400).send("No uid provided");
  }

  try {
    const files = req.files as Express.Multer.File[];

    const uploadPromises = files.map((file) => {
      const imagePath = file.path;
      const fileName = file.filename;
      return uploadMain(imagePath, fileName, uid);
    });

    await Promise.all(uploadPromises);

    const s3Keys = files.map((file) => `${uid}/${file.filename}`);

    res.status(200).json({
      keys: s3Keys,
    });
  } catch (error) {
    console.error("Error uploading images to S3:", error);
    for (const file of req.files as Express.Multer.File[]) {
      try {
        await unlink(file.path);
      } catch (cleanupError) {
        console.error("Failed to cleanup file:", file.path, cleanupError);
      }
    }
    res.status(500).send("Error uploading images to S3");
  }
});

app.get("/get", async (req, res) => {
  console.log("GET request received");
  try {
    const date = req.query.date;

    if (typeof date !== "string" || !date) {
      res.status(400).send("Error: missing parameters");
      return;
    }

    const result = await s3.send(
      new GetObjectCommand({
        Bucket: process.env.BUCKET_NAME,
        Key: date,
      })
    );

    if (!result.Body) {
      res.status(500).send("Error: No response from S3");
      return;
    }

    const readableObj = result.Body as Readable;
    res.setHeader("Content-Type", result.ContentType || "image/png");
    res.setHeader("Content-Length", result.ContentLength as number);

    readableObj.pipe(res);

    readableObj.on("end", () => {
      console.log("Readable stream ended");
    });

    readableObj.on("error", (err) => {
      console.error("Error in readable stream:", err);
      res.status(500).send("Error in readable stream");
    });
  } catch (error) {
    console.error("Error processing request:", error);
    res.status(500).send("Error processing request");
  }
});

export default app;
