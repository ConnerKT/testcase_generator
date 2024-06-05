const mongoose = require('mongoose');
const Problem = require('./models/oldProblem'); 
const dotenv = require('dotenv').config();
const data = require('./allProblems.json');

if (!process.env.DB_TESTING) {
    console.error("DB_URI_TESTING is not defined in the environment variables.");
    process.exit(1);
  }

mongoose.connect(process.env.DB_TESTING, { useNewUrlParser: true, useUnifiedTopology: true });

async function seedProblems() {
  try {
    for (let i = 0; i < data.length; i++) {
      const problemData = data[i];
      
      const newProblem = new Problem({
        id: problemData.questionId,
        title: problemData.questionTitle,
        difficulty: problemData.difficulty,
        description: problemData.question,
        link: problemData.link
      });

      await newProblem.save();
      console.log(`Problem ${i + 1} seeded successfully`);
    }

    console.log('All problems seeded successfully');
  } catch (err) {
    console.error('Error seeding problems:', err);
  } finally {
    mongoose.connection.close();
  }
}

seedProblems();