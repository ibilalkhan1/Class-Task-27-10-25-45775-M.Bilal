import express from "express";
import dotenv from "dotenv";
import mysql from "mysql2/promise";
import bcrypt from "bcrypt";
import jwt from "jsonwebtoken";

dotenv.config();
const app = express();
app.use(express.json());

const pool = await mysql.createPool({
  host: process.env.MYSQL_HOST || "mysql",
  user: process.env.MYSQL_USER || "root",
  password: process.env.MYSQL_PASSWORD || "rootpass",
  database: process.env.MYSQL_DB || "usersdb",
  waitForConnections: true,
  connectionLimit: 10
});

await pool.query(`CREATE TABLE IF NOT EXISTS users(
  id INT AUTO_INCREMENT PRIMARY KEY,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  name VARCHAR(255) DEFAULT NULL
)`);

app.get("/health", (_req, res) => res.json({ ok: true }));

app.post("/signup", async (req, res) => {
  const { email, password, name } = req.body || {};
  if (!email || !password) return res.status(400).json({ error: "email+password required" });
  const hash = await bcrypt.hash(password, 10);
  try {
    await pool.query("INSERT INTO users(email,password_hash,name) VALUES (?,?,?)", [email, hash, name || null]);
    res.status(201).json({ message: "user created" });
  } catch {
    res.status(400).json({ error: "email exists?" });
  }
});

app.post("/login", async (req, res) => {
  const { email, password } = req.body || {};
  const [rows] = await pool.query("SELECT * FROM users WHERE email=?", [email]);
  const u = rows[0];
  if (!u) return res.status(401).json({ error: "invalid" });
  const ok = await bcrypt.compare(password, u.password_hash);
  if (!ok) return res.status(401).json({ error: "invalid" });
  const token = jwt.sign({ sub: u.id, email: u.email }, process.env.JWT_SECRET || "devjwt", { expiresIn: "1h" });
  res.json({ token });
});

app.put("/profile", async (req, res) => {
  const { token, name } = req.body || {};
  try {
    const p = jwt.verify(token, process.env.JWT_SECRET || "devjwt");
    await pool.query("UPDATE users SET name=? WHERE id=?", [name || null, p.sub]);
    res.json({ message: "updated" });
  } catch {
    res.status(401).json({ error: "bad token" });
  }
});

app.post("/password-reset", (req, res) => {
  res.json({ message: `reset link would be sent to ${req.body?.email}` });
});

export default app;
