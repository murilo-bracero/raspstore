const conn = new Mongo(process.env.DB_URI);

const db = conn.getDB('raspstore');

const collectionExists = db.getCollectionNames().includes('users');

if(!collectionExists){
    db.createCollection('users');
}

db.getCollection('users').createIndex({"username": 1}, {unique: true});
db.getCollection('users').createIndex({"user_id": 1}, {unique: true});