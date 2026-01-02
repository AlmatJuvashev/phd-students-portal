CREATE TABLE lti_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    kid VARCHAR(255) NOT NULL UNIQUE, -- Key ID exposed in JWKS
    private_key TEXT NOT NULL, -- PEM encoded RSA Private Key
    public_key TEXT NOT NULL, -- PEM encoded RSA Public Key
    algorithm VARCHAR(50) DEFAULT 'RS256',
    use VARCHAR(50) DEFAULT 'sig',
    
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    expires_at TIMESTAMP WITH TIME ZONE -- Optional key rotation
);

CREATE INDEX idx_lti_keys_kid ON lti_keys(kid);
