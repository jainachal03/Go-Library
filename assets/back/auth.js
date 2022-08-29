const express = require('express')

const { MongoClient, ServerApiVersion, Double } = require('mongodb');
const mongoose = require("mongoose")

const dotenv = require('dotenv');

dotenv.config()

const bodyParser = require('body-parser');

const app = express();
app.use(bodyParser.urlencoded({ extended: true }));

const uri = `mongodb+srv://achal:${process.env.PASSWORD}@cluster0.byopwfw.mongodb.net/?retryWrites=true&w=majority`;
const client = new MongoClient(uri, { useNewUrlParser: true, useUnifiedTopology: true, serverApi: ServerApiVersion.v1 });
client.connect(err => {
});
db =client.db("bookHandler")

app.get("/", (req, res)=>{
    res.send("HelloWorld")
})
app.post('/auth', (req, res)=>{
    getAllDocs = async (query) => {
        return await db.collection(coll).find(query).toArray()
    }
    const body = req.body
    let username = body['username']
    let  password = body['pwd']
    const query = {"username":username, "password":password}
    const response = getAllDocs(query)
    if(response){
        res.sendStatus(200);
    }else{
        res.sendStatus(400);
    }
    return

})

app.listen(3000)

