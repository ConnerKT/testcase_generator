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

function transformFunctionSignatures(problem) {
  const signatures = problem.functionSignatures;
  const transformed = [];

  for (const [language, func] of Object.entries(signatures)) {
    // Extract function name from the function string
    let name;
    if (language.toLowerCase() === "python") {
      name = func.match(/def (\w+)/)[1];
    } else if (language.toLowerCase() === "javascript") {
      name = func.match(/function (\w+)/)[1];
    } else {
      // Add more languages and their corresponding regex patterns here if needed
      name = "unknown";
    }

    transformed.push({
      name: name,
      language: language,
      value: func,
    });
  }

  return transformed;
}

// Example usage
const problem = {
  "functionSignatures": {
    "python": "def intToRoman(num: int) -> str:\n\n\n    pass",
    "javascript": "function intToRoman(num) {\n\n\n}"
  }
};

async function seeding() {
  try {
  await mongoose.connect(process.env.DB_URI, {
    useNewUrlParser: true,
    useUnifiedTopology: true,
  });
  console.log('Connected to MongoDB');

  // Fetch all documents
  const allDocuments = await Problem.find({});
  
  for (let i = 0; i < allDocuments.length; i++) {
    const problem = allDocuments[i];
    const transformedSignatures = transformFunctionSignatures(problem);

    // Update the document's functionSignatures field
    problem.functionSignatures = transformedSignatures;
    await problem.save();
    
    console.log(`Updated document with id: ${problem._id}`);
  }
    // Close the database connection
    await mongoose.connection.close();
    console.log('Connection closed');
  } catch (err) {
    console.error('Error connecting to MongoDB or fetching problems:', err);
  }
}

seeding();


// Make a function that takes in a problem object and returns an array formatted like this:

// "functionSignatures": [
//   {
//     "name": "function name",
//     "language": "the language",
//     "value": "the function"
//   }
// ]

// For example:

// I will pass in this object

// "functionSignatures": {
//   "python": "def intToRoman(num: int) -> str:\n\n\n    pass",
//   "javascript": "function intToRoman(num) {\n\n\n}"
// }

// and the function will return

// "functionSignatures": [
//   {
//     "name": "intToRoman",
//     "language": "Python",
//     "value": "def intToRoman(num: int) -> str:\n\n\n    pass"
//   },
//   {
//     "name": "intToRoman",
//     "language": "JavaScript",
//     "value": "function intToRoman(num) {\n\n\n}"
//   }
// ]