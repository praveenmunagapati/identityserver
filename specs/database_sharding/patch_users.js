var i = 0;
db = new Mongo().getDB("itsyouonline-idserver-db");
db.users.find().forEach(function(doc) {
    db.users.update(
        { "_id": doc._id },
        { "$set": { "country": i % 2 === 0 ? "EU" : "RU" } }
    );
    i++
});
