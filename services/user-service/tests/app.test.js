import request from "supertest";
import app from "../src/app.js";
describe("user-service", ()=>{
  it("health", async ()=>{
    const r = await request(app).get("/health");
    expect(r.statusCode).toBe(200);
    expect(r.body.ok).toBe(true);
  });
});
