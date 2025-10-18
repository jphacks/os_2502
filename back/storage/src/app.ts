import express from "express";
import cookieParser from "cookie-parser";
import logger from "morgan";
import cors from "cors";

import collageRouters from "./app/root";

require("dotenv").config();
const app = express();
const router = express.Router();
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
app.use(router);

app.use("/upload", collageRouters);

// 接続確認用
console.log("Hello World!!!");
app.get("/ping", (_, req) => {
  req.json({ message: "connected" });
});

export default app;
