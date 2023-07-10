use identity

db.users.stats()
db.users.createIndex({ email: 1 });
db.users.createIndex({ '$**': 'text' });
db.users.getIndexes();

db.blacklists.stats()
db.blacklists.createIndex({ access_token: 1 });
db.blacklists.createIndex({ '$**': 'text' });
db.blacklists.getIndexes();

db.user_groups.stats()
db.user_groups.createIndex({ creator_id: 1 });
db.user_groups.createIndex({ '$**': 'text' });
db.user_groups.getIndexes();

db.memberships.stats()
db.memberships.createIndex({ group_id: 1 });
db.memberships.createIndex({ user_id: 1 });
db.memberships.createIndex({ '$**': 'text' });
db.memberships.getIndexes();

db.user_memberships.stats()
db.user_memberships.createIndex({ membership_id: 1 });
db.user_memberships.createIndex({ user_id: 1 });
db.user_memberships.createIndex({ group_id: 1 });
db.user_memberships.createIndex({ '$**': 'text' });
db.user_memberships.getIndexes();

db.group_memberships.stats()
db.group_memberships.createIndex({ membership_id: 1 });
db.group_memberships.createIndex({ group_id: 1 });
db.group_memberships.createIndex({ user_id: 1 });
db.group_memberships.createIndex({ '$**': 'text' });
db.group_memberships.getIndexes();
