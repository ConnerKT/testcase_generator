const mongoose = require('mongoose');


const problemSchema = new mongoose.Schema({
    id: Number,
    title: String,
    difficulty: String,
    description: String,
    link: String,
},{
    versionKey: false 
}, {
    collection: 'challenges' // Specify the collection name here
  });

const Problem = mongoose.model('Problem', problemSchema);

module.exports = Problem;