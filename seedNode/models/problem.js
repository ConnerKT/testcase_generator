const mongoose = require('mongoose');

const functionSignatureSchema = new mongoose.Schema({
  name: { type: String, required: true },
  language: { type: String, required: true },
  value: { type: String, required: true }
}, { _id: false });

const testCaseSchema = new mongoose.Schema({
  id: { type: String, required: true },
  input: {
    nums: { type: [Number], required: true },
    target: { type: Number, required: true }
  },
  output: { type: [Number], required: true }
}, { _id: false });

const problemSchema = new mongoose.Schema({
  title: { type: String, required: true },
  difficulty: { type: String, required: true },
  description: { type: String, required: true },
  link: { type: String, required: true },
  functionSignatures: [functionSignatureSchema],
  testCases: [testCaseSchema]
});

const Problem = mongoose.model('Problem', problemSchema);

module.exports = Problem;
