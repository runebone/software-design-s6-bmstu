// Meh, hardcode
db = connect("mongodb://mongo_user:password@user-mongo:27017");

db.users.insertMany([
  { id: "aaaaaaaa-dddd-0000-0000-000000000000", username: "admin", email: "admin@gmail.com", role: "admin", password_hash: "$2a$10$tMXCVXRe/SHD0TzRkO107.ezmuNaDPrdLZpb4u6zOQbwbha2wRY3S" },
]);
