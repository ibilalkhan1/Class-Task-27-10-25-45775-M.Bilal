import request from "supertest";
import app from "../src/app.js";

const email = "ciuser@example.com";
const pass  = "Pass123!";
let token;

describe("user-service auth flow", () => {
  it("signup 201 or 400 when exists", async () => {
    const r = await request(app).post("/signup").send({email, password: pass, name: "CI User"});
    expect([201,400]).toContain(r.statusCode);
  });
  it("login returns token", async () => {
    const r = await request(app).post("/login").send({email, password: pass});
    expect(r.statusCode).toBe(200);
    expect(r.body.token).toBeTruthy();
    token = r.body.token;
  });
});
