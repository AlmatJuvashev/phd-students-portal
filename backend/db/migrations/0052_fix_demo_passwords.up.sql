-- Fix passwords for all demo users to 'demopassword123!'
-- Generated hash: $2a$10$vA.YkvqwsCK8u12BnYjBnOsx/8TVxPHTi101GSQOC4ZERlV2vZiwK
UPDATE users 
SET password_hash = '$2a$10$vA.YkvqwsCK8u12BnYjBnOsx/8TVxPHTi101GSQOC4ZERlV2vZiwK'
WHERE username LIKE 'demo.%' OR username LIKE 'dr.%';
