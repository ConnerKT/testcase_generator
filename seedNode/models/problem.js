const mongoose = require('mongoose');

const problemSchema = new mongoose.Schema({
  title: {
    type: String,
    required: true
  },
  difficulty: {
    type: String,
    required: true
  },
  description: {
    type: String,
    required: true
  },
  link: {
    type: String,
    required: true
  },
  functionSignatures: {
      python: String,
      javascript: String
  }
});

const Problem = mongoose.model('challenges', problemSchema);

module.exports = Problem;
