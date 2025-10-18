import { load } from "ts-dotenv";
import { GetObjectCommand, S3Client } from "@aws-sdk/client-s3";
import { Upload } from "@aws-sdk/lib-storage";
import fs from "fs";

const env = load({
  BUCKET_NAME: String,
  // AWS_PROFILE: String, // 削除
  S3_ENDPOINT: String,
  AWS_ACCESS_KEY_ID: String, // 追加
  AWS_SECRET_ACCESS_KEY: String, // 追加
  AWS_REGION: { type: String, default: "ap-northeast-1" }, // 追加
});

const s3Client = new S3Client({
  region: env.AWS_REGION,
  credentials: {
    accessKeyId: env.AWS_ACCESS_KEY_ID,
    secretAccessKey: env.AWS_SECRET_ACCESS_KEY,
  },
  endpoint: env.S3_ENDPOINT,
  forcePathStyle: true,
});

const uploadObject = async (filePath: string, key: string) => {
  const fileStream = fs.readFileSync(filePath);

  try {
    const inputParams = new Upload({
      client: s3Client,
      params: {
        Bucket: env.BUCKET_NAME,
        Key: key,
        Body: fileStream,
      },
    });

    const result = await inputParams.done();
    console.log("アップロード完了", result.Location);
    return result.Location;
  } catch (error) {
    console.error("アップロード失敗", error);
  }
};

const uploadMain = async (filePath: string, fileName: string, uid: string) => {
  try {
    const s3Key = `${uid}/${fileName}`;
    return await uploadObject(filePath, s3Key);
  } catch (error) {
    console.error("Error uploading image to S3:", error);
    throw error;
  }
};

const downLoadObject = async (key: string) => {
  try {
    const getObjectParams = new GetObjectCommand({
      Bucket: env.BUCKET_NAME,
      Key: key,
    });
    const result = await s3Client.send(getObjectParams);
    return result;
  } catch (error) {
    console.error("ダウンロード失敗", error);
    throw error; // エラーを再スローして呼び出し元にエラーを伝える
  }
};

// ダウンロード処理のメイン関数
const downloadMain = async (filePath: string) => {
  try {
    console.log("Starting download process");
    return await downLoadObject(filePath);
  } catch (error) {
    console.error("Error during download process:", error);
    throw error; // エラーを再スローして呼び出し元にエラーを伝える
  }
};

export { uploadMain, downloadMain };
