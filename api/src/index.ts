import "module-alias/register";
import "dotenv/config";
import "@/utils/timers.util";
import express from "express";
import { createServer } from "http";
import { Server } from "socket.io";
import logger from "@/loggers/default.logger";

const app = express();
const httpServer = createServer(app);
const io = new Server(httpServer, {
  /* options */
});

io.on("connection", (socket) => {
  // ...
});

httpServer.listen(3000, () => {
  logger.info(`app chạy lúc ${new Date().toLocaleString()}`);
});
