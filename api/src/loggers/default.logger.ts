import pino from "pino";
// import rfs from "rotating-file-stream";
import path from "path";
import fs from "fs";
import { timer } from "@/utils/timers.util";
const rfs = require("rotating-file-stream");

const logDirectory = path.join(__dirname, "..", "..", "logs");
if (!fs.existsSync(logDirectory)) {
  fs.mkdirSync(logDirectory);
}

// Tạo stream ghi log theo ngày
const stream = rfs.createStream(
  (time: string | number | Date, index: any) => {
    let date = timer();
    if (time) {
      date = timer(time);
    }
    return `app-${date.format("YYYYMMDD")}.log`;
  },
  {
    interval: "1d", // Xoay log mỗi ngày
    path: logDirectory,
    compress: "gzip", // Nén các file log cũ
    maxFiles: 14, // Giữ lại các file log trong 14 ngày
  }
);

const logger = pino(
  {
    level: "info",
    timestamp: pino.stdTimeFunctions.isoTime,
  },
  stream
);

export default logger;
