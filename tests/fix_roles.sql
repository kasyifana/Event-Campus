-- SQL Commands to Fix User Roles for Testing

-- 1. Fix Admin Account
-- Sets role to 'admin', approves account, and resets password to 'admin123456'
-- Password hash for 'admin123456' is: $2a$10$YourHashedPasswordHere (Placeholder - use the one below if you can generate it, otherwise just update role)
-- Since we can't easily generate bcrypt hash here, we will assume you might need to register a new admin if you don't know the password.
-- BUT, if you just registered it via API and know the password, just run:
UPDATE users 
SET role = 'admin', is_approved = true 
WHERE email = 'admin@eventcampus.com';

-- 2. Fix Organizer Account
-- Sets role to 'organisasi' and approves account
UPDATE users 
SET role = 'organisasi', is_approved = true 
WHERE email = 'organizer@eventcampus.com';

-- 3. Verify
SELECT email, role, is_approved FROM users WHERE email IN ('admin@eventcampus.com', 'organizer@eventcampus.com');
