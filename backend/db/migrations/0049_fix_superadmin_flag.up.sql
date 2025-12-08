-- Fix: Ensure superadmin users have is_superadmin flag set to true
UPDATE users SET is_superadmin = true WHERE role = 'superadmin';
