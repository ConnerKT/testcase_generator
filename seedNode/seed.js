require("dotenv").config();

const express = require("express");
const mongoose = require("mongoose");
const app = express();
const port = 3000;
const Problem = require("./models/problem");

app.get("/", (req, res) => {
  res.send("Hello World");
});

app.listen(port, () => {
  console.log("My app is running on port", port);
});

async function seeding() {
  try {
    await mongoose.connect(process.env.DB_URI, {
      useNewUrlParser: true,
      useUnifiedTopology: true,
    });
    console.log("Connected to MongoDB");

    const allProblems = await Problem.find({}, 'functionSignatures');
    console.log(allProblems);

    await mongoose.connection.close();
    console.log("Connection closed");


  } catch (err) {
    console.error("Error connecting to MongoDB or fetching problems:", err);
  }
}

seeding();
