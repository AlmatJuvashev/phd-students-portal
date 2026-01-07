# Comprehensive Guide to Enterprise Features

**Document Version:** 1.0  
**Created:** January 6, 2026  
**Purpose:** Detailed explanation of missing and partial enterprise features from Competitive Analysis

---

## Table of Contents

1. [SSO/SAML/OAuth2 - Enterprise Authentication](#1-ssosaml-oauth2---enterprise-authentication)
2. [WebSocket/Real-time Communication](#2-websocketreal-time-communication) *(coming soon)*
3. [Video Conferencing Integration](#3-video-conferencing-integration) *(coming soon)*
4. [LTI 1.3 (Learning Tools Interoperability)](#4-lti-13-learning-tools-interoperability) *(coming soon)*
5. [SCORM Support](#5-scorm-support) *(coming soon)*
6. [xAPI/Learning Record Store](#6-xapilearning-record-store) *(coming soon)*
7. [Mobile Applications & Push Notifications](#7-mobile-applications--push-notifications) *(coming soon)*
8. [WCAG 2.1 AA Accessibility](#8-wcag-21-aa-accessibility) *(coming soon)*
9. [Gamification System](#9-gamification-system) *(coming soon)*
10. [AI Tutoring & Adaptive Learning](#10-ai-tutoring--adaptive-learning) *(coming soon)*

---

## 1. SSO/SAML/OAuth2 - Enterprise Authentication

### 1.1 Definition

**Single Sign-On (SSO)** — это технология аутентификации, которая позволяет пользователям входить в несколько связанных, но независимых программных систем, используя единый набор учетных данных (логин и пароль). После однократной аутентификации пользователь получает доступ ко всем подключенным системам без необходимости повторного ввода пароля.

**SAML (Security Assertion Markup Language)** — это открытый стандарт на основе XML для обмена данными аутентификации и авторизации между сторонами, в частности между поставщиком удостоверений (Identity Provider, IdP) и поставщиком услуг (Service Provider, SP).

**OAuth 2.0** — это протокол авторизации, который позволяет приложениям получать ограниченный доступ к учетным записям пользователей на сторонних сервисах (Google, Microsoft, Facebook и др.).

**OIDC (OpenID Connect)** — это слой идентификации, построенный поверх OAuth 2.0, который добавляет возможности аутентификации пользователей.

#### Ключевые компоненты:

| Компонент | Описание | Пример |
|-----------|----------|--------|
| **Identity Provider (IdP)** | Система, хранящая и проверяющая удостоверения пользователей | Okta, Azure AD, Google Workspace |
| **Service Provider (SP)** | Приложение, требующее аутентификации (наша платформа) | PhD Student Portal |
| **Assertion** | XML-документ с данными об аутентифицированном пользователе | SAML Response |
| **Token** | JWT или другой токен для OAuth/OIDC | Access Token, ID Token |

#### Поток аутентификации SAML:

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│              │     │              │     │              │
│    User      │     │   Service    │     │   Identity   │
│   Browser    │     │   Provider   │     │   Provider   │
│              │     │  (Our App)   │     │  (Okta/AD)   │
└──────┬───────┘     └──────┬───────┘     └──────┬───────┘
       │                    │                    │
       │  1. Access App     │                    │
       │───────────────────>│                    │
       │                    │                    │
       │  2. Redirect to IdP│                    │
       │<───────────────────│                    │
       │                    │                    │
       │  3. Login Request  │                    │
       │────────────────────────────────────────>│
       │                    │                    │
       │  4. User Authenticates (login/password) │
       │<────────────────────────────────────────│
       │                    │                    │
       │  5. SAML Assertion │                    │
       │────────────────────────────────────────>│
       │                    │                    │
       │  6. Redirect with Assertion             │
       │<────────────────────────────────────────│
       │                    │                    │
       │  7. Submit Assertion                    │
       │───────────────────>│                    │
       │                    │                    │
       │  8. Validate & Create Session           │
       │                    │───────────────────>│
       │                    │<───────────────────│
       │                    │                    │
       │  9. Access Granted │                    │
       │<───────────────────│                    │
       │                    │                    │
```

---

### 1.2 Почему важно реализовать

#### Бизнес-причины:

| Причина | Описание | Влияние |
|---------|----------|---------|
| **Deal Breaker** | Без SSO университеты НЕ будут покупать систему | Критическое — потеря 100% enterprise-клиентов |
| **Compliance** | Требования безопасности корпоративных политик | Юридические и регуляторные требования |
| **IT Policy** | Централизованное управление доступом обязательно | Без SSO невозможна интеграция в IT-инфраструктуру |
| **Security Audit** | Проверки безопасности требуют централизованной аутентификации | Провал аудита = отказ от внедрения |

#### Технические причины:

1. **Централизованное управление пользователями** — IT-отдел управляет всеми пользователями в одном месте (Active Directory, Okta)
2. **Единая политика паролей** — требования к сложности, ротации, MFA применяются централизованно
3. **Автоматический provisioning/deprovisioning** — при увольнении сотрудника доступ отключается мгновенно во всех системах
4. **Audit trail** — централизованное логирование всех входов для compliance

#### Статистика рынка:

```
📊 Данные исследований:
• 92% организаций с >1000 сотрудников используют SSO
• 78% университетов требуют SAML для новых систем
• Средняя экономия: $1.2M/год на helpdesk (сброс паролей)
• 50% сокращение security incidents при использовании SSO
```

---

### 1.3 Что дает конечному пользователю

#### Для студентов и преподавателей:

| Преимущество | Описание | Пользовательский опыт |
|--------------|----------|----------------------|
| **Один пароль** | Не нужно запоминать отдельный пароль для LMS | Снижение когнитивной нагрузки |
| **Быстрый вход** | Один клик — и вы в системе | Экономия времени (5-10 сек на каждый вход) |
| **Бесшовный опыт** | Переход между системами без повторного входа | Productivity boost |
| **Безопасность** | Защита корпоративным IdP с MFA | Меньше риск взлома аккаунта |
| **Self-service** | Сброс пароля через единый портал | Не нужно обращаться в поддержку |

#### Для IT-администраторов:

| Преимущество | Описание | Экономия |
|--------------|----------|----------|
| **Централизация** | Управление всеми пользователями в одном месте | 60% сокращение времени на user management |
| **Onboarding** | Автоматическое создание аккаунтов при добавлении в AD | С часов до минут |
| **Offboarding** | Мгновенное отключение доступа при увольнении | Security compliance |
| **Отчетность** | Единый audit log для всех систем | Упрощение compliance-проверок |
| **Helpdesk** | Меньше тикетов "забыл пароль" | 30-50% снижение нагрузки |

#### Конкретные сценарии использования:

```
Сценарий 1: Студент первого курса
├─ Без SSO: Получает email, создает аккаунт, придумывает пароль,
│           проходит верификацию, забывает пароль, сбрасывает...
└─ С SSO:   Входит через университетский портал одним кликом ✓

Сценарий 2: Преподаватель с 5 системами
├─ Без SSO: 5 разных паролей, постоянные сбросы, risk of reuse
└─ С SSO:   Один вход утром = доступ везде весь день ✓

Сценарий 3: Увольнение сотрудника
├─ Без SSO: HR уведомляет IT, IT вручную отключает в каждой системе,
│           риск пропустить систему, security vulnerability
└─ С SSO:   Деактивация в AD = мгновенное отключение везде ✓
```

---

### 1.4 Какие пользователи предпочитают эту функцию

#### Целевые сегменты:

| Сегмент | Размер | Требование SSO | Приоритет |
|---------|--------|----------------|-----------|
| **Университеты** | >5,000 студентов | Обязательно | P0 |
| **Корпорации** | >500 сотрудников | Обязательно | P0 |
| **Государственные учреждения** | Любой размер | Обязательно (compliance) | P0 |
| **Школы K-12** | >1,000 учеников | Google SSO обязательно | P0 |
| **Малый бизнес** | <100 сотрудников | Желательно | P2 |
| **Индивидуальные пользователи** | 1 | Не требуется | P3 |

#### Типичные требования по отраслям:

**Высшее образование (Университеты):**
- SAML 2.0 с Shibboleth/InCommon Federation
- Интеграция с существующим IdP (Azure AD, Okta, PingFederate)
- Атрибуты: eduPersonAffiliation, eduPersonPrincipalName
- MFA через университетский IdP

**Корпоративное обучение:**
- Azure AD / Okta / OneLogin интеграция
- SCIM для автоматического provisioning
- Just-in-time (JIT) user provisioning
- Группы и роли из IdP

**K-12 образование:**
- Google Workspace for Education SSO (приоритет)
- Clever SSO для школьных систем
- ClassLink integration
- Простой интерфейс для детей

**Государственные учреждения (Казахстан/СНГ):**
- Интеграция с национальными ID-системами (ЭЦП)
- ГОСТ-совместимое шифрование
- Локальное хранение данных
- Аудит всех действий

---

### 1.5 Как интегрировать в наше приложение

#### Архитектура решения:

```
┌─────────────────────────────────────────────────────────────────┐
│                    Authentication Layer                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────────────┐  │
│  │   SAML 2.0   │  │   OAuth2/    │  │   LDAP/Active       │  │
│  │   Handler    │  │   OIDC       │  │   Directory         │  │
│  └──────┬───────┘  └──────┬───────┘  └──────────┬───────────┘  │
│         │                 │                      │              │
│         └────────────────┬┴──────────────────────┘              │
│                          │                                       │
│                  ┌───────▼────────┐                             │
│                  │  Identity      │                             │
│                  │  Service       │                             │
│                  └───────┬────────┘                             │
│                          │                                       │
│         ┌────────────────┼────────────────┐                     │
│         │                │                │                     │
│  ┌──────▼──────┐  ┌──────▼──────┐  ┌─────▼─────┐              │
│  │   User      │  │  Session    │  │  Audit    │              │
│  │   Linking   │  │  Manager    │  │  Logger   │              │
│  └─────────────┘  └─────────────┘  └───────────┘              │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

#### Требуемые изменения в базе данных:

```sql
-- Таблица провайдеров идентификации
CREATE TABLE identity_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    
    -- Основные настройки
    provider_type VARCHAR(20) NOT NULL, -- 'saml', 'oauth2', 'ldap', 'oidc'
    name VARCHAR(100) NOT NULL,
    display_name VARCHAR(200),
    
    -- Конфигурация (зависит от типа)
    config JSONB NOT NULL,
    /*
    Для SAML:
    {
        "entity_id": "https://idp.university.edu",
        "sso_url": "https://idp.university.edu/sso",
        "slo_url": "https://idp.university.edu/slo",
        "certificate": "-----BEGIN CERTIFICATE-----...",
        "attribute_mapping": {
            "email": "urn:oid:0.9.2342.19200300.100.1.3",
            "first_name": "urn:oid:2.5.4.42",
            "last_name": "urn:oid:2.5.4.4"
        }
    }
    
    Для OAuth2/OIDC:
    {
        "client_id": "abc123",
        "client_secret": "encrypted:...",
        "auth_url": "https://accounts.google.com/o/oauth2/auth",
        "token_url": "https://oauth2.googleapis.com/token",
        "userinfo_url": "https://openidconnect.googleapis.com/v1/userinfo",
        "scopes": ["openid", "email", "profile"]
    }
    
    Для LDAP:
    {
        "host": "ldap.university.edu",
        "port": 636,
        "use_ssl": true,
        "bind_dn": "cn=service,dc=university,dc=edu",
        "bind_password": "encrypted:...",
        "base_dn": "ou=users,dc=university,dc=edu",
        "user_filter": "(uid=%s)"
    }
    */
    
    -- Настройки поведения
    is_active BOOLEAN DEFAULT true,
    is_default BOOLEAN DEFAULT false,
    allow_password_login BOOLEAN DEFAULT true, -- Разрешить вход по паролю
    auto_create_users BOOLEAN DEFAULT true,    -- JIT provisioning
    auto_update_profile BOOLEAN DEFAULT true,  -- Обновлять профиль при входе
    
    -- Сопоставление ролей
    role_mapping JSONB DEFAULT '{}',
    /*
    {
        "admin_group": "admin",
        "teacher_group": "instructor",
        "student_group": "student"
    }
    */
    
    -- Метаданные
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    
    UNIQUE(tenant_id, name)
);

-- Связь внешних идентификаторов с пользователями
CREATE TABLE external_identities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider_id UUID NOT NULL REFERENCES identity_providers(id) ON DELETE CASCADE,
    
    -- Внешний идентификатор
    external_id VARCHAR(255) NOT NULL, -- Subject/NameID из IdP
    email VARCHAR(255),
    
    -- Дополнительные данные от IdP
    metadata JSONB DEFAULT '{}',
    /*
    {
        "groups": ["students", "cs-department"],
        "attributes": {
            "department": "Computer Science",
            "student_id": "2024001234"
        }
    }
    */
    
    -- Временные метки
    linked_at TIMESTAMP DEFAULT NOW(),
    last_login_at TIMESTAMP,
    
    UNIQUE(provider_id, external_id)
);

-- Токены MFA (если IdP не предоставляет)
CREATE TABLE mfa_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Тип MFA
    mfa_type VARCHAR(20) NOT NULL, -- 'totp', 'sms', 'email', 'webauthn'
    
    -- Секрет (зашифрованный)
    secret_encrypted TEXT NOT NULL,
    
    -- Для TOTP
    backup_codes TEXT[], -- Зашифрованные резервные коды
    
    -- Для WebAuthn
    credential_id BYTEA,
    public_key BYTEA,
    
    -- Статус
    is_verified BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    
    -- Метаданные
    created_at TIMESTAMP DEFAULT NOW(),
    last_used_at TIMESTAMP,
    
    UNIQUE(user_id, mfa_type)
);

-- Индексы для производительности
CREATE INDEX idx_identity_providers_tenant ON identity_providers(tenant_id);
CREATE INDEX idx_external_identities_user ON external_identities(user_id);
CREATE INDEX idx_external_identities_external_id ON external_identities(external_id);
CREATE INDEX idx_mfa_tokens_user ON mfa_tokens(user_id);
```

#### Реализация на Go (основные компоненты):

```go
// internal/auth/sso/types.go
package sso

import (
    "context"
    "time"
)

// SSOProvider определяет интерфейс для всех SSO провайдеров
type SSOProvider interface {
    // GetAuthURL возвращает URL для перенаправления пользователя
    GetAuthURL(state string, redirectURL string) (string, error)
    
    // HandleCallback обрабатывает callback от IdP
    HandleCallback(ctx context.Context, code string) (*UserInfo, error)
    
    // ValidateToken проверяет токен (для API-вызовов)
    ValidateToken(ctx context.Context, token string) (*UserInfo, error)
    
    // GetProviderInfo возвращает метаданные провайдера
    GetProviderInfo() ProviderInfo
}

// UserInfo содержит информацию о пользователе от IdP
type UserInfo struct {
    ExternalID    string            `json:"external_id"`
    Email         string            `json:"email"`
    EmailVerified bool              `json:"email_verified"`
    FirstName     string            `json:"first_name"`
    LastName      string            `json:"last_name"`
    DisplayName   string            `json:"display_name"`
    Groups        []string          `json:"groups"`
    Attributes    map[string]string `json:"attributes"`
    RawData       map[string]any    `json:"raw_data"`
}

// ProviderInfo содержит метаданные провайдера
type ProviderInfo struct {
    Type        string `json:"type"`
    Name        string `json:"name"`
    DisplayName string `json:"display_name"`
    IconURL     string `json:"icon_url"`
}
```

```go
// internal/auth/sso/saml_provider.go
package sso

import (
    "context"
    "crypto/x509"
    "encoding/base64"
    "fmt"
    
    "github.com/crewjam/saml"
    "github.com/crewjam/saml/samlsp"
)

type SAMLProvider struct {
    sp              *samlsp.Middleware
    config          SAMLConfig
    attributeMapping map[string]string
}

type SAMLConfig struct {
    EntityID        string `json:"entity_id"`
    SSOURL          string `json:"sso_url"`
    SLOURL          string `json:"slo_url"`
    Certificate     string `json:"certificate"`
    AttributeMapping map[string]string `json:"attribute_mapping"`
}

func NewSAMLProvider(config SAMLConfig, spEntityID string, acsURL string) (*SAMLProvider, error) {
    // Парсим сертификат IdP
    certData, err := base64.StdEncoding.DecodeString(config.Certificate)
    if err != nil {
        return nil, fmt.Errorf("invalid certificate: %w", err)
    }
    
    cert, err := x509.ParseCertificate(certData)
    if err != nil {
        return nil, fmt.Errorf("failed to parse certificate: %w", err)
    }
    
    // Создаем IdP descriptor
    idpMetadata := &saml.EntityDescriptor{
        EntityID: config.EntityID,
        IDPSSODescriptors: []saml.IDPSSODescriptor{
            {
                SingleSignOnServices: []saml.Endpoint{
                    {
                        Binding:  saml.HTTPRedirectBinding,
                        Location: config.SSOURL,
                    },
                },
                KeyDescriptors: []saml.KeyDescriptor{
                    {
                        Use: "signing",
                        KeyInfo: saml.KeyInfo{
                            Certificate: base64.StdEncoding.EncodeToString(cert.Raw),
                        },
                    },
                },
            },
        },
    }
    
    // Настройки SP остаются в middleware
    // ...
    
    return &SAMLProvider{
        config:          config,
        attributeMapping: config.AttributeMapping,
    }, nil
}

func (p *SAMLProvider) HandleCallback(ctx context.Context, samlResponse string) (*UserInfo, error) {
    // Декодируем и валидируем SAML Response
    // Извлекаем атрибуты согласно маппингу
    // Возвращаем UserInfo
    
    return &UserInfo{
        // Заполняем из SAML Assertion
    }, nil
}
```

```go
// internal/auth/sso/oauth_provider.go
package sso

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    
    "golang.org/x/oauth2"
)

type OAuthProvider struct {
    config       *oauth2.Config
    userInfoURL  string
    providerType string
}

type OAuthConfig struct {
    ClientID     string   `json:"client_id"`
    ClientSecret string   `json:"client_secret"`
    AuthURL      string   `json:"auth_url"`
    TokenURL     string   `json:"token_url"`
    UserInfoURL  string   `json:"userinfo_url"`
    Scopes       []string `json:"scopes"`
}

func NewGoogleOAuthProvider(clientID, clientSecret, redirectURL string) *OAuthProvider {
    return &OAuthProvider{
        config: &oauth2.Config{
            ClientID:     clientID,
            ClientSecret: clientSecret,
            RedirectURL:  redirectURL,
            Scopes:       []string{"openid", "email", "profile"},
            Endpoint: oauth2.Endpoint{
                AuthURL:  "https://accounts.google.com/o/oauth2/auth",
                TokenURL: "https://oauth2.googleapis.com/token",
            },
        },
        userInfoURL:  "https://openidconnect.googleapis.com/v1/userinfo",
        providerType: "google",
    }
}

func NewMicrosoftOAuthProvider(clientID, clientSecret, tenantID, redirectURL string) *OAuthProvider {
    return &OAuthProvider{
        config: &oauth2.Config{
            ClientID:     clientID,
            ClientSecret: clientSecret,
            RedirectURL:  redirectURL,
            Scopes:       []string{"openid", "email", "profile", "User.Read"},
            Endpoint: oauth2.Endpoint{
                AuthURL:  fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/authorize", tenantID),
                TokenURL: fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID),
            },
        },
        userInfoURL:  "https://graph.microsoft.com/v1.0/me",
        providerType: "microsoft",
    }
}

func (p *OAuthProvider) GetAuthURL(state string, redirectURL string) (string, error) {
    if redirectURL != "" {
        p.config.RedirectURL = redirectURL
    }
    return p.config.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

func (p *OAuthProvider) HandleCallback(ctx context.Context, code string) (*UserInfo, error) {
    // Обмениваем код на токен
    token, err := p.config.Exchange(ctx, code)
    if err != nil {
        return nil, fmt.Errorf("failed to exchange code: %w", err)
    }
    
    // Получаем информацию о пользователе
    client := p.config.Client(ctx, token)
    resp, err := client.Get(p.userInfoURL)
    if err != nil {
        return nil, fmt.Errorf("failed to get user info: %w", err)
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("user info request failed: %d", resp.StatusCode)
    }
    
    var rawData map[string]any
    if err := json.NewDecoder(resp.Body).Decode(&rawData); err != nil {
        return nil, fmt.Errorf("failed to decode user info: %w", err)
    }
    
    // Маппинг полей зависит от провайдера
    return p.mapUserInfo(rawData), nil
}

func (p *OAuthProvider) mapUserInfo(data map[string]any) *UserInfo {
    userInfo := &UserInfo{
        RawData: data,
    }
    
    switch p.providerType {
    case "google":
        userInfo.ExternalID, _ = data["sub"].(string)
        userInfo.Email, _ = data["email"].(string)
        userInfo.EmailVerified, _ = data["email_verified"].(bool)
        userInfo.FirstName, _ = data["given_name"].(string)
        userInfo.LastName, _ = data["family_name"].(string)
        userInfo.DisplayName, _ = data["name"].(string)
        
    case "microsoft":
        userInfo.ExternalID, _ = data["id"].(string)
        userInfo.Email, _ = data["mail"].(string)
        if userInfo.Email == "" {
            userInfo.Email, _ = data["userPrincipalName"].(string)
        }
        userInfo.FirstName, _ = data["givenName"].(string)
        userInfo.LastName, _ = data["surname"].(string)
        userInfo.DisplayName, _ = data["displayName"].(string)
    }
    
    return userInfo
}
```

```go
// internal/auth/sso/service.go
package sso

import (
    "context"
    "fmt"
    "time"
    
    "github.com/google/uuid"
)

type SSOService struct {
    repo       SSORepository
    userRepo   UserRepository
    providers  map[uuid.UUID]SSOProvider
}

func NewSSOService(repo SSORepository, userRepo UserRepository) *SSOService {
    return &SSOService{
        repo:      repo,
        userRepo:  userRepo,
        providers: make(map[uuid.UUID]SSOProvider),
    }
}

// AuthenticateSSO обрабатывает SSO callback и возвращает/создает пользователя
func (s *SSOService) AuthenticateSSO(
    ctx context.Context, 
    tenantID uuid.UUID,
    providerID uuid.UUID, 
    userInfo *UserInfo,
) (*User, error) {
    // 1. Ищем существующую связку external_identity
    externalIdentity, err := s.repo.FindExternalIdentity(ctx, providerID, userInfo.ExternalID)
    if err != nil && !IsNotFoundError(err) {
        return nil, fmt.Errorf("failed to find external identity: %w", err)
    }
    
    // 2. Если связка найдена — возвращаем пользователя
    if externalIdentity != nil {
        user, err := s.userRepo.FindByID(ctx, externalIdentity.UserID)
        if err != nil {
            return nil, fmt.Errorf("failed to find user: %w", err)
        }
        
        // Обновляем время последнего входа
        if err := s.repo.UpdateLastLogin(ctx, externalIdentity.ID); err != nil {
            // Логируем, но не прерываем
        }
        
        // Обновляем профиль если настроено
        provider, err := s.repo.FindProviderByID(ctx, providerID)
        if err == nil && provider.AutoUpdateProfile {
            s.updateUserProfile(ctx, user, userInfo)
        }
        
        return user, nil
    }
    
    // 3. Связки нет — ищем пользователя по email
    user, err := s.userRepo.FindByEmail(ctx, tenantID, userInfo.Email)
    if err != nil && !IsNotFoundError(err) {
        return nil, fmt.Errorf("failed to find user by email: %w", err)
    }
    
    // 4. Пользователь найден по email — связываем
    if user != nil {
        if err := s.linkExternalIdentity(ctx, user.ID, providerID, userInfo); err != nil {
            return nil, fmt.Errorf("failed to link identity: %w", err)
        }
        return user, nil
    }
    
    // 5. Пользователя нет — проверяем JIT provisioning
    provider, err := s.repo.FindProviderByID(ctx, providerID)
    if err != nil {
        return nil, fmt.Errorf("failed to find provider: %w", err)
    }
    
    if !provider.AutoCreateUsers {
        return nil, fmt.Errorf("user not found and auto-creation disabled")
    }
    
    // 6. Создаем нового пользователя (JIT provisioning)
    user, err = s.createUserFromSSO(ctx, tenantID, provider, userInfo)
    if err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    // 7. Связываем с external identity
    if err := s.linkExternalIdentity(ctx, user.ID, providerID, userInfo); err != nil {
        // Откатываем создание пользователя
        s.userRepo.Delete(ctx, user.ID)
        return nil, fmt.Errorf("failed to link new user identity: %w", err)
    }
    
    return user, nil
}

func (s *SSOService) createUserFromSSO(
    ctx context.Context,
    tenantID uuid.UUID,
    provider *IdentityProvider,
    userInfo *UserInfo,
) (*User, error) {
    // Определяем роль на основе маппинга групп
    role := s.determineRole(provider.RoleMapping, userInfo.Groups)
    
    user := &User{
        ID:        uuid.New(),
        TenantID:  tenantID,
        Email:     userInfo.Email,
        FirstName: userInfo.FirstName,
        LastName:  userInfo.LastName,
        Role:      role,
        IsActive:  true,
        CreatedAt: time.Now(),
    }
    
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    
    return user, nil
}

func (s *SSOService) determineRole(roleMapping map[string]string, groups []string) string {
    for group, role := range roleMapping {
        for _, userGroup := range groups {
            if userGroup == group {
                return role
            }
        }
    }
    return "student" // Default role
}
```

#### API Endpoints:

```go
// internal/handlers/sso_handler.go
package handlers

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type SSOHandler struct {
    ssoService *sso.SSOService
}

// GET /api/v1/auth/sso/providers
// Возвращает список доступных SSO провайдеров для tenant
func (h *SSOHandler) ListProviders(c *gin.Context) {
    tenantID := c.GetString("tenant_id")
    
    providers, err := h.ssoService.ListActiveProviders(c.Request.Context(), uuid.MustParse(tenantID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    // Возвращаем только публичную информацию
    response := make([]gin.H, len(providers))
    for i, p := range providers {
        response[i] = gin.H{
            "id":           p.ID,
            "name":         p.Name,
            "display_name": p.DisplayName,
            "type":         p.ProviderType,
            "icon_url":     getProviderIcon(p.ProviderType),
        }
    }
    
    c.JSON(http.StatusOK, response)
}

// GET /api/v1/auth/sso/:provider_id/login
// Инициирует SSO flow, перенаправляет на IdP
func (h *SSOHandler) InitiateLogin(c *gin.Context) {
    providerID := c.Param("provider_id")
    redirectURL := c.Query("redirect_url")
    
    // Генерируем state для CSRF защиты
    state := generateSecureState()
    
    // Сохраняем state в Redis с TTL
    h.ssoService.SaveState(c.Request.Context(), state, redirectURL)
    
    // Получаем URL для перенаправления
    authURL, err := h.ssoService.GetAuthURL(c.Request.Context(), uuid.MustParse(providerID), state)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// POST /api/v1/auth/sso/:provider_id/callback
// Обрабатывает callback от IdP
func (h *SSOHandler) HandleCallback(c *gin.Context) {
    providerID := c.Param("provider_id")
    
    // Для SAML — SAMLResponse в body
    // Для OAuth — code в query
    
    var userInfo *sso.UserInfo
    var err error
    
    provider, _ := h.ssoService.GetProvider(c.Request.Context(), uuid.MustParse(providerID))
    
    switch provider.ProviderType {
    case "saml":
        samlResponse := c.PostForm("SAMLResponse")
        userInfo, err = h.ssoService.HandleSAMLCallback(c.Request.Context(), uuid.MustParse(providerID), samlResponse)
    case "oauth2", "oidc":
        code := c.Query("code")
        state := c.Query("state")
        
        // Проверяем state
        if !h.ssoService.ValidateState(c.Request.Context(), state) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
            return
        }
        
        userInfo, err = h.ssoService.HandleOAuthCallback(c.Request.Context(), uuid.MustParse(providerID), code)
    }
    
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    
    // Аутентифицируем/создаем пользователя
    tenantID := c.GetString("tenant_id")
    user, err := h.ssoService.AuthenticateSSO(c.Request.Context(), uuid.MustParse(tenantID), uuid.MustParse(providerID), userInfo)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    
    // Создаем сессию
    token, err := h.sessionService.CreateSession(c.Request.Context(), user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    // Перенаправляем на frontend с токеном
    redirectURL := h.ssoService.GetSavedRedirectURL(c.Request.Context(), state)
    c.Redirect(http.StatusTemporaryRedirect, redirectURL+"?token="+token)
}
```

#### Конфигурация для Admin UI:

```typescript
// frontend/src/types/sso.ts
export interface IdentityProvider {
  id: string;
  tenantId: string;
  providerType: 'saml' | 'oauth2' | 'oidc' | 'ldap';
  name: string;
  displayName: string;
  isActive: boolean;
  isDefault: boolean;
  allowPasswordLogin: boolean;
  autoCreateUsers: boolean;
  autoUpdateProfile: boolean;
  config: SAMLConfig | OAuthConfig | LDAPConfig;
  roleMapping: Record<string, string>;
}

export interface SAMLConfig {
  entityId: string;
  ssoUrl: string;
  sloUrl?: string;
  certificate: string;
  attributeMapping: {
    email: string;
    firstName: string;
    lastName: string;
    groups?: string;
  };
}

export interface OAuthConfig {
  clientId: string;
  clientSecret: string;
  authUrl: string;
  tokenUrl: string;
  userinfoUrl: string;
  scopes: string[];
}
```

---

### 1.6 Сложность освоения и реализации

#### Оценка сложности:

| Аспект | Сложность | Комментарий |
|--------|-----------|-------------|
| **SAML 2.0** | 🔴 Высокая | Сложный XML, криптография, много edge cases |
| **OAuth 2.0** | 🟡 Средняя | Хорошо документирован, много библиотек |
| **OIDC** | 🟡 Средняя | Надстройка над OAuth, стандартизированы claims |
| **LDAP** | 🟡 Средняя | Legacy протокол, нужно понимать схемы |
| **MFA** | 🟢 Низкая | Много готовых библиотек (TOTP, WebAuthn) |

#### Временные оценки:

```
Реализация SAML 2.0 SP:
├── Изучение стандарта: 3-5 дней
├── Базовая реализация: 5-7 дней
├── Тестирование с IdP: 3-5 дней
├── Edge cases & debugging: 3-5 дней
└── Итого: 2-3 недели

Реализация OAuth2/OIDC:
├── Google OAuth: 1-2 дня
├── Microsoft OAuth: 1-2 дня
├── Generic OIDC: 2-3 дня
└── Итого: 1 неделя

Реализация LDAP:
├── Базовое подключение: 2-3 дня
├── Синхронизация групп: 2-3 дня
├── Тестирование с AD: 2-3 дня
└── Итого: 1 неделя

Admin UI для настройки:
├── Формы конфигурации: 3-5 дней
├── Тестирование соединения: 2-3 дня
└── Итого: 1 неделя

Общее время: 5-6 недель (один разработчик)
```

#### Типичные проблемы и решения:

| Проблема | Причина | Решение |
|----------|---------|---------|
| "Invalid signature" в SAML | Неправильный сертификат или namespace | Проверить XML canonicalization, сертификат |
| OAuth redirect mismatch | Несоответствие redirect_uri | Точное совпадение URL в настройках |
| LDAP timeout | Firewall или неверный порт | Проверить сетевое подключение, 389/636 |
| User not found after SSO | JIT provisioning отключен | Включить auto_create_users |
| Роль не назначается | Неверный маппинг групп | Проверить названия групп в IdP |

#### Требуемые навыки разработчика:

- ✅ Понимание криптографии (подписи, сертификаты X.509)
- ✅ Опыт работы с XML (для SAML)
- ✅ Знание HTTP security (CORS, cookies, CSRF)
- ✅ Понимание федеративной идентификации
- ✅ Опыт отладки сетевых протоколов

---

### 1.7 Источники для дальнейшего изучения

#### Официальные спецификации:

| Ресурс | Ссылка | Описание |
|--------|--------|----------|
| **SAML 2.0 Specification** | [docs.oasis-open.org/security/saml](https://docs.oasis-open.org/security/saml/v2.0/) | Полная спецификация SAML |
| **OAuth 2.0 RFC 6749** | [tools.ietf.org/html/rfc6749](https://tools.ietf.org/html/rfc6749) | Базовый стандарт OAuth 2.0 |
| **OpenID Connect** | [openid.net/specs](https://openid.net/specs/openid-connect-core-1_0.html) | Спецификация OIDC |
| **LDAP RFC 4511** | [tools.ietf.org/html/rfc4511](https://tools.ietf.org/html/rfc4511) | Протокол LDAP |

#### Библиотеки и инструменты (Go):

| Библиотека | Ссылка | Назначение |
|------------|--------|------------|
| **crewjam/saml** | [github.com/crewjam/saml](https://github.com/crewjam/saml) | SAML 2.0 для Go |
| **golang.org/x/oauth2** | [pkg.go.dev/golang.org/x/oauth2](https://pkg.go.dev/golang.org/x/oauth2) | OAuth 2.0 client |
| **coreos/go-oidc** | [github.com/coreos/go-oidc](https://github.com/coreos/go-oidc) | OIDC для Go |
| **go-ldap/ldap** | [github.com/go-ldap/ldap](https://github.com/go-ldap/ldap) | LDAP client для Go |
| **pquerna/otp** | [github.com/pquerna/otp](https://github.com/pquerna/otp) | TOTP/HOTP для MFA |

#### Документация IdP провайдеров:

| Провайдер | Документация | Комментарий |
|-----------|--------------|-------------|
| **Okta** | [developer.okta.com](https://developer.okta.com/docs/) | Отличная документация для интеграции |
| **Azure AD** | [docs.microsoft.com/azure/active-directory](https://docs.microsoft.com/en-us/azure/active-directory/) | SAML и OAuth для Microsoft |
| **Google Workspace** | [developers.google.com/identity](https://developers.google.com/identity) | OAuth и OIDC для Google |
| **Auth0** | [auth0.com/docs](https://auth0.com/docs/) | Universal login, много примеров |
| **Keycloak** | [keycloak.org/documentation](https://www.keycloak.org/documentation) | Open source IdP |

#### Обучающие ресурсы:

| Ресурс | Ссылка | Формат |
|--------|--------|--------|
| **SAML for Web Developers** | [samltool.com](https://www.samltool.com/generic_sso_req.php) | Интерактивные инструменты |
| **OAuth.net** | [oauth.net/2/](https://oauth.net/2/) | Гайды и best practices |
| **Auth0 Blog** | [auth0.com/blog](https://auth0.com/blog/) | Статьи о безопасности |
| **Keycloak Tutorials** | [youtube.com/@Keycloak](https://www.youtube.com/@Keycloak) | Видео-туториалы |

#### Инструменты для тестирования:

| Инструмент | Назначение | Ссылка |
|------------|------------|--------|
| **SAML Tracer** | Browser extension для отладки SAML | Chrome/Firefox extension |
| **jwt.io** | Декодирование и валидация JWT | [jwt.io](https://jwt.io/) |
| **Postman** | Тестирование OAuth flows | [postman.com](https://www.postman.com/) |
| **SimpleSAMLphp** | Test IdP для разработки | [simplesamlphp.org](https://simplesamlphp.org/) |
| **Keycloak** | Локальный IdP для тестов | Docker: `quay.io/keycloak/keycloak` |

---

### 1.8 Чек-лист реализации

```
Phase 1: SAML 2.0 (Week 1-2)
□ Изучить SAML 2.0 спецификацию
□ Настроить тестовый IdP (Keycloak/SimpleSAML)
□ Реализовать SP metadata endpoint
□ Реализовать ACS (Assertion Consumer Service)
□ Реализовать SLO (Single Logout)
□ Тестирование с Okta/Azure AD
□ Документация для администраторов

Phase 2: OAuth2/OIDC (Week 2-3)
□ Реализовать Google OAuth
□ Реализовать Microsoft OAuth
□ Реализовать generic OIDC provider
□ State management для CSRF
□ Token refresh handling
□ Тестирование всех провайдеров

Phase 3: LDAP (Week 4)
□ Реализовать LDAP bind/search
□ Поддержка SSL/TLS
□ Синхронизация групп
□ Тестирование с Active Directory

Phase 4: Admin UI (Week 4-5)
□ Форма добавления SAML IdP
□ Форма добавления OAuth provider
□ Форма добавления LDAP connection
□ Test connection функционал
□ Role mapping UI
□ Документация для пользователей

Phase 5: MFA (Week 5-6)
□ TOTP (Google Authenticator)
□ Backup codes
□ Recovery flow
□ Admin-enforced MFA policy
```

---

## 2. LDAP Integration (Lightweight Directory Access Protocol)

### 2.1 Определение

**LDAP (Lightweight Directory Access Protocol)** — это открытый, кроссплатформенный протокол прикладного уровня для доступа и управления распределёнными службами каталогов через IP-сеть. LDAP является стандартом де-факто для корпоративных систем управления идентификацией и хранения учетных данных пользователей.

**Active Directory (AD)** — это реализация службы каталогов от Microsoft, которая использует LDAP как один из основных протоколов доступа. AD является стандартом в корпоративной среде Windows.

**Directory Service** — это иерархическая база данных, оптимизированная для операций чтения, которая хранит информацию о пользователях, группах, компьютерах и других объектах сети.

#### Ключевые концепции LDAP:

| Термин | Описание | Пример |
|--------|----------|--------|
| **DN (Distinguished Name)** | Уникальный идентификатор объекта в каталоге | `cn=John Doe,ou=Users,dc=university,dc=edu` |
| **Base DN** | Корневая точка для поиска в каталоге | `dc=university,dc=edu` |
| **Bind DN** | Учетная запись для подключения к LDAP | `cn=service,ou=Apps,dc=university,dc=edu` |
| **OU (Organizational Unit)** | Контейнер для группировки объектов | `ou=Students`, `ou=Faculty` |
| **Attribute** | Свойство объекта (email, имя, группы) | `mail`, `givenName`, `memberOf` |
| **Filter** | Выражение для поиска объектов | `(&(objectClass=user)(mail=*@university.edu))` |

#### Структура LDAP-каталога:

```
dc=university,dc=edu (Domain)
├── ou=Users
│   ├── ou=Students
│   │   ├── cn=John Doe
│   │   │   ├── uid: jdoe
│   │   │   ├── mail: jdoe@university.edu
│   │   │   ├── givenName: John
│   │   │   ├── sn: Doe
│   │   │   └── memberOf: cn=CS-Students,ou=Groups,...
│   │   └── cn=Jane Smith
│   └── ou=Faculty
│       └── cn=Dr. Brown
├── ou=Groups
│   ├── cn=Administrators
│   ├── cn=CS-Students
│   └── cn=PhD-Candidates
└── ou=Applications
    └── cn=LMS-Service
```

#### Поток аутентификации LDAP:

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│              │     │              │     │              │
│    User      │     │  Our App     │     │   LDAP       │
│   Browser    │     │  (Backend)   │     │   Server     │
│              │     │              │     │  (AD/OpenLDAP)│
└──────┬───────┘     └──────┬───────┘     └──────┬───────┘
       │                    │                    │
       │ 1. Login Form      │                    │
       │   (username/pass)  │                    │
       │───────────────────>│                    │
       │                    │                    │
       │                    │ 2. Bind as Service │
       │                    │    Account         │
       │                    │───────────────────>│
       │                    │                    │
       │                    │ 3. Bind OK         │
       │                    │<───────────────────│
       │                    │                    │
       │                    │ 4. Search User     │
       │                    │   (by username)    │
       │                    │───────────────────>│
       │                    │                    │
       │                    │ 5. User DN Found   │
       │                    │<───────────────────│
       │                    │                    │
       │                    │ 6. Bind as User    │
       │                    │   (verify password)│
       │                    │───────────────────>│
       │                    │                    │
       │                    │ 7. Bind Success    │
       │                    │<───────────────────│
       │                    │                    │
       │                    │ 8. Get User Attrs  │
       │                    │   (groups, email)  │
       │                    │───────────────────>│
       │                    │                    │
       │                    │ 9. User Attributes │
       │                    │<───────────────────│
       │                    │                    │
       │ 10. Session Token  │                    │
       │<───────────────────│                    │
       │                    │                    │
```

---

### 2.2 Почему важно реализовать

#### Бизнес-причины:

| Причина | Описание | Влияние на бизнес |
|---------|----------|-------------------|
| **Enterprise Standard** | 95% крупных организаций используют AD/LDAP | Без LDAP — нет enterprise-клиентов |
| **Централизованное управление** | IT-отдел управляет пользователями в одном месте | Снижение TCO (Total Cost of Ownership) |
| **Compliance** | Требование корпоративных политик безопасности | Прохождение security audit |
| **Legacy Integration** | Многие организации не готовы к OAuth/SAML | Расширение рынка |

#### Технические причины:

1. **Единый источник правды (Single Source of Truth)** — все пользователи, группы и права в одном каталоге
2. **Автоматическая синхронизация** — при добавлении пользователя в AD он автоматически появляется в LMS
3. **Групповые политики** — назначение ролей на основе членства в группах AD
4. **Offboarding** — при отключении пользователя в AD доступ в LMS блокируется автоматически
5. **Password Policy** — политика паролей управляется централизованно в AD

#### Статистика и факты:

```
📊 Рыночные данные:
• 95% организаций Fortune 500 используют Active Directory
• 78% университетов используют LDAP для аутентификации
• Средний enterprise имеет 15+ приложений, интегрированных с AD
• LDAP-интеграция сокращает helpdesk-тикеты на 40%
• ROI интеграции: 200-300% за первый год (экономия на управлении пользователями)
```

#### Сравнение с SSO/OAuth:

| Аспект | LDAP | SSO/OAuth |
|--------|------|-----------|
| **Протокол** | Binary over TCP | HTTP/HTTPS |
| **Возраст** | 1993+ (legacy) | 2012+ (modern) |
| **Использование** | На убывание | На рост |
| **Firewall** | Требует открытия портов 389/636 | Работает через 443 |
| **Сложность** | Средняя | Низкая-Средняя |
| **Когда нужен** | Legacy системы, локальные сети | Облачные системы |

**Важно:** Многие организации используют LDAP и SSO параллельно. LDAP — для внутренних систем, SSO — для облачных.

---

### 2.3 Что дает конечному пользователю

#### Для студентов и преподавателей:

| Преимущество | Описание | Пользовательский опыт |
|--------------|----------|----------------------|
| **Один пароль** | Тот же пароль, что и для входа в компьютер | Не нужно запоминать новый пароль |
| **Автоматическая регистрация** | Аккаунт создается автоматически | Нет процесса регистрации |
| **Актуальные данные** | ФИО, email, факультет из AD | Профиль всегда актуален |
| **Автоматические роли** | Роль назначается по группе в AD | Правильные права с первого входа |
| **Смена пароля** | Смена в AD = смена везде | Один раз сменил — везде обновилось |

#### Для IT-администраторов:

| Преимущество | Описание | Экономия |
|--------------|----------|----------|
| **Нет дублирования** | Пользователи не создаются вручную в LMS | 90% экономия времени на onboarding |
| **Централизованный контроль** | Отключил в AD = отключил везде | Безопасность при увольнении |
| **Групповое управление** | Добавил в группу AD = получил роль в LMS | Массовое назначение прав |
| **Аудит** | Все входы логируются в AD | Единый audit trail |
| **Password Reset** | Пользователи сбрасывают пароль в AD | Меньше тикетов в helpdesk |

#### Конкретные сценарии:

```
Сценарий 1: Новый студент
├─ Без LDAP: HR создает в AD → IT создает в LMS → студент получает 2 пароля
└─ С LDAP:   HR создает в AD → студент входит в LMS с тем же паролем ✓

Сценарий 2: Преподаватель стал деканом
├─ Без LDAP: IT меняет роль в AD → IT вручную меняет роль в LMS
└─ С LDAP:   IT добавляет в группу Deans в AD → роль в LMS меняется автоматически ✓

Сценарий 3: Студент отчислен
├─ Без LDAP: HR отключает в AD → IT забывает отключить в LMS → security risk
└─ С LDAP:   HR отключает в AD → при попытке входа LMS отклоняет автоматически ✓

Сценарий 4: Смена пароля
├─ Без LDAP: Пользователь меняет пароль в AD, потом в LMS (часто забывает)
└─ С LDAP:   Пользователь меняет пароль в AD → работает везде ✓
```

---

### 2.4 Какие пользователи предпочитают эту функцию

#### Целевые сегменты:

| Сегмент | Размер | Требование LDAP | Приоритет |
|---------|--------|-----------------|-----------|
| **Крупные университеты** | >10,000 студентов | Обязательно | P0 |
| **Государственные учреждения** | Любой | Обязательно (compliance) | P0 |
| **Корпорации** | >500 сотрудников | Часто обязательно | P0 |
| **Средние университеты** | 1,000-10,000 | Желательно | P1 |
| **K-12 школы** | Разный | Иногда (если есть AD) | P2 |
| **Малые организации** | <100 | Редко | P3 |

#### Типичные требования по отраслям:

**Высшее образование (Университеты):**
```
• Active Directory или OpenLDAP
• Интеграция с существующим IdM (Identity Management)
• Синхронизация групп: студенты по факультетам, преподаватели по кафедрам
• Атрибуты: студенческий ID, email, факультет, год поступления
• Требование: пользователь не должен существовать только в LMS
```

**Корпоративное обучение:**
```
• Active Directory (99% случаев)
• Интеграция с HR-системами через AD
• Группы: по отделам, по уровню доступа, по проектам
• Атрибуты: employee ID, department, manager, job title
• Требование: соответствие политикам безопасности
```

**Государственные учреждения:**
```
• LDAP с ГОСТ-совместимым шифрованием (для СНГ)
• Строгий audit trail всех операций
• Интеграция с национальными системами идентификации
• Требование: сертификация по стандартам безопасности
```

#### Вопросы при продаже:

```
Типичные вопросы от enterprise-клиентов:
1. "Поддерживаете ли вы интеграцию с Active Directory?"
2. "Можем ли мы синхронизировать группы из AD?"
3. "Поддерживается ли LDAPS (LDAP over SSL)?"
4. "Можно ли использовать AD для аутентификации, но не создавать пользователей автоматически?"
5. "Как часто происходит синхронизация атрибутов?"

Без LDAP ответ на первый вопрос = "Нет" = потеря сделки
```

---

### 2.5 Как интегрировать в наше приложение

#### Архитектура решения:

```
┌─────────────────────────────────────────────────────────────────┐
│                     LDAP Integration Layer                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────────┐  ┌──────────────────┐  ┌───────────────┐ │
│  │  LDAP Connection │  │   User Search    │  │    Group      │ │
│  │     Pool         │  │   & Bind         │  │    Sync       │ │
│  └────────┬─────────┘  └────────┬─────────┘  └───────┬───────┘ │
│           │                     │                    │          │
│           └─────────────────────┼────────────────────┘          │
│                                 │                               │
│                    ┌────────────▼────────────┐                  │
│                    │      LDAP Service       │                  │
│                    │  (Authentication +      │                  │
│                    │   Attribute Mapping)    │                  │
│                    └────────────┬────────────┘                  │
│                                 │                               │
│         ┌───────────────────────┼───────────────────────┐       │
│         │                       │                       │       │
│  ┌──────▼──────┐    ┌───────────▼───────────┐  ┌───────▼─────┐ │
│  │   Auth      │    │    User Provisioning  │  │   Role      │ │
│  │   Handler   │    │    (JIT / Scheduled)  │  │   Mapper    │ │
│  └─────────────┘    └───────────────────────┘  └─────────────┘ │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

#### Схема базы данных:

```sql
-- Конфигурация LDAP-подключений
CREATE TABLE ldap_configurations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    
    -- Основные настройки подключения
    name VARCHAR(100) NOT NULL,
    host VARCHAR(255) NOT NULL,           -- ldap.university.edu
    port INTEGER DEFAULT 389,              -- 389 для LDAP, 636 для LDAPS
    use_ssl BOOLEAN DEFAULT false,         -- LDAPS
    use_starttls BOOLEAN DEFAULT false,    -- STARTTLS
    
    -- Учетные данные для bind
    bind_dn VARCHAR(500) NOT NULL,         -- cn=service,ou=Apps,dc=university,dc=edu
    bind_password_encrypted TEXT NOT NULL, -- Зашифрованный пароль
    
    -- Настройки поиска
    base_dn VARCHAR(500) NOT NULL,         -- dc=university,dc=edu
    user_search_base VARCHAR(500),         -- ou=Users,dc=university,dc=edu
    user_search_filter VARCHAR(500) DEFAULT '(uid=%s)',
    group_search_base VARCHAR(500),        -- ou=Groups,dc=university,dc=edu
    group_search_filter VARCHAR(500) DEFAULT '(objectClass=group)',
    
    -- Маппинг атрибутов
    attribute_mapping JSONB NOT NULL DEFAULT '{
        "username": "uid",
        "email": "mail",
        "first_name": "givenName",
        "last_name": "sn",
        "display_name": "displayName",
        "employee_id": "employeeID",
        "department": "department",
        "groups": "memberOf"
    }',
    
    -- Маппинг групп на роли
    role_mapping JSONB DEFAULT '{
        "cn=Admins,ou=Groups,dc=university,dc=edu": "admin",
        "cn=Faculty,ou=Groups,dc=university,dc=edu": "instructor",
        "cn=Students,ou=Groups,dc=university,dc=edu": "student"
    }',
    
    -- Настройки поведения
    is_active BOOLEAN DEFAULT true,
    allow_password_login BOOLEAN DEFAULT true, -- Разрешить вход по паролю (не только LDAP)
    auto_create_users BOOLEAN DEFAULT true,    -- JIT provisioning
    auto_update_profile BOOLEAN DEFAULT true,  -- Обновлять профиль при входе
    sync_groups BOOLEAN DEFAULT true,          -- Синхронизировать группы
    
    -- Настройки синхронизации
    sync_enabled BOOLEAN DEFAULT false,        -- Периодическая синхронизация
    sync_interval_minutes INTEGER DEFAULT 60,  -- Интервал синхронизации
    last_sync_at TIMESTAMP,
    last_sync_status VARCHAR(50),
    last_sync_error TEXT,
    
    -- Таймауты
    connection_timeout_seconds INTEGER DEFAULT 10,
    search_timeout_seconds INTEGER DEFAULT 30,
    
    -- Метаданные
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    created_by UUID REFERENCES users(id)
);

-- Связь LDAP-пользователей с локальными
CREATE TABLE ldap_user_mappings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ldap_config_id UUID NOT NULL REFERENCES ldap_configurations(id) ON DELETE CASCADE,
    
    -- LDAP идентификаторы
    ldap_dn VARCHAR(500) NOT NULL,           -- cn=John Doe,ou=Users,dc=university,dc=edu
    ldap_uid VARCHAR(255) NOT NULL,          -- jdoe
    ldap_object_guid VARCHAR(100),           -- Уникальный ID в AD (для отслеживания переименований)
    
    -- Кэшированные атрибуты
    cached_attributes JSONB DEFAULT '{}',
    cached_groups TEXT[],                    -- Список DN групп
    
    -- Синхронизация
    last_sync_at TIMESTAMP,
    last_login_at TIMESTAMP,
    
    UNIQUE(ldap_config_id, ldap_dn),
    UNIQUE(ldap_config_id, ldap_uid)
);

-- Лог синхронизации
CREATE TABLE ldap_sync_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ldap_config_id UUID NOT NULL REFERENCES ldap_configurations(id),
    
    sync_type VARCHAR(20) NOT NULL,          -- 'full', 'incremental', 'user'
    started_at TIMESTAMP NOT NULL,
    completed_at TIMESTAMP,
    status VARCHAR(20) NOT NULL,             -- 'running', 'success', 'failed'
    
    -- Статистика
    users_found INTEGER DEFAULT 0,
    users_created INTEGER DEFAULT 0,
    users_updated INTEGER DEFAULT 0,
    users_disabled INTEGER DEFAULT 0,
    errors_count INTEGER DEFAULT 0,
    
    -- Детали
    error_details JSONB,
    
    created_at TIMESTAMP DEFAULT NOW()
);

-- Индексы
CREATE INDEX idx_ldap_config_tenant ON ldap_configurations(tenant_id);
CREATE INDEX idx_ldap_mapping_user ON ldap_user_mappings(user_id);
CREATE INDEX idx_ldap_mapping_uid ON ldap_user_mappings(ldap_config_id, ldap_uid);
CREATE INDEX idx_ldap_sync_logs_config ON ldap_sync_logs(ldap_config_id);
```

#### Реализация на Go:

```go
// internal/auth/ldap/client.go
package ldap

import (
    "context"
    "crypto/tls"
    "fmt"
    "strings"
    "time"

    "github.com/go-ldap/ldap/v3"
)

// LDAPClient управляет подключением к LDAP-серверу
type LDAPClient struct {
    config     *LDAPConfig
    conn       *ldap.Conn
    connected  bool
}

type LDAPConfig struct {
    Host                string
    Port                int
    UseSSL              bool
    UseStartTLS         bool
    BindDN              string
    BindPassword        string
    BaseDN              string
    UserSearchBase      string
    UserSearchFilter    string   // "(uid=%s)" или "(sAMAccountName=%s)"
    GroupSearchBase     string
    GroupSearchFilter   string
    AttributeMapping    map[string]string
    RoleMapping         map[string]string
    ConnectionTimeout   time.Duration
    SearchTimeout       time.Duration
}

func NewLDAPClient(config *LDAPConfig) *LDAPClient {
    return &LDAPClient{
        config: config,
    }
}

// Connect устанавливает соединение с LDAP-сервером
func (c *LDAPClient) Connect() error {
    var conn *ldap.Conn
    var err error

    address := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)

    if c.config.UseSSL {
        // LDAPS (порт 636)
        tlsConfig := &tls.Config{
            ServerName: c.config.Host,
            MinVersion: tls.VersionTLS12,
        }
        conn, err = ldap.DialTLS("tcp", address, tlsConfig)
    } else {
        // Обычный LDAP (порт 389)
        conn, err = ldap.Dial("tcp", address)
        
        if err == nil && c.config.UseStartTLS {
            // Обновление до TLS через STARTTLS
            err = conn.StartTLS(&tls.Config{
                ServerName: c.config.Host,
                MinVersion: tls.VersionTLS12,
            })
        }
    }

    if err != nil {
        return fmt.Errorf("failed to connect to LDAP: %w", err)
    }

    // Bind с сервисной учетной записью
    err = conn.Bind(c.config.BindDN, c.config.BindPassword)
    if err != nil {
        conn.Close()
        return fmt.Errorf("failed to bind to LDAP: %w", err)
    }

    c.conn = conn
    c.connected = true
    return nil
}

// Close закрывает соединение
func (c *LDAPClient) Close() {
    if c.conn != nil {
        c.conn.Close()
        c.connected = false
    }
}

// Authenticate проверяет учетные данные пользователя
func (c *LDAPClient) Authenticate(ctx context.Context, username, password string) (*LDAPUser, error) {
    if !c.connected {
        if err := c.Connect(); err != nil {
            return nil, err
        }
    }

    // 1. Поиск пользователя по username
    user, err := c.SearchUser(ctx, username)
    if err != nil {
        return nil, fmt.Errorf("user not found: %w", err)
    }

    // 2. Попытка bind с найденным DN и паролем пользователя
    userConn, err := c.createConnection()
    if err != nil {
        return nil, err
    }
    defer userConn.Close()

    err = userConn.Bind(user.DN, password)
    if err != nil {
        return nil, fmt.Errorf("invalid credentials: %w", err)
    }

    // 3. Получение групп пользователя
    groups, err := c.GetUserGroups(ctx, user.DN)
    if err != nil {
        // Логируем, но не прерываем
        groups = []string{}
    }
    user.Groups = groups

    return user, nil
}

// SearchUser ищет пользователя по username
func (c *LDAPClient) SearchUser(ctx context.Context, username string) (*LDAPUser, error) {
    // Формируем фильтр поиска
    filter := fmt.Sprintf(c.config.UserSearchFilter, ldap.EscapeFilter(username))
    
    searchBase := c.config.UserSearchBase
    if searchBase == "" {
        searchBase = c.config.BaseDN
    }

    // Атрибуты для получения
    attributes := []string{"dn"}
    for _, attr := range c.config.AttributeMapping {
        attributes = append(attributes, attr)
    }

    searchRequest := ldap.NewSearchRequest(
        searchBase,
        ldap.ScopeWholeSubtree,
        ldap.NeverDerefAliases,
        0,    // Size limit
        int(c.config.SearchTimeout.Seconds()),
        false,
        filter,
        attributes,
        nil,
    )

    result, err := c.conn.Search(searchRequest)
    if err != nil {
        return nil, fmt.Errorf("LDAP search failed: %w", err)
    }

    if len(result.Entries) == 0 {
        return nil, fmt.Errorf("user not found: %s", username)
    }

    if len(result.Entries) > 1 {
        return nil, fmt.Errorf("multiple users found for: %s", username)
    }

    entry := result.Entries[0]
    return c.entryToUser(entry), nil
}

// GetUserGroups получает список групп пользователя
func (c *LDAPClient) GetUserGroups(ctx context.Context, userDN string) ([]string, error) {
    // Для Active Directory используем memberOf атрибут
    // Для OpenLDAP может потребоваться другой подход
    
    searchRequest := ldap.NewSearchRequest(
        userDN,
        ldap.ScopeBaseObject,
        ldap.NeverDerefAliases,
        0,
        int(c.config.SearchTimeout.Seconds()),
        false,
        "(objectClass=*)",
        []string{"memberOf"},
        nil,
    )

    result, err := c.conn.Search(searchRequest)
    if err != nil {
        return nil, err
    }

    if len(result.Entries) == 0 {
        return []string{}, nil
    }

    return result.Entries[0].GetAttributeValues("memberOf"), nil
}

// SyncAllUsers синхронизирует всех пользователей
func (c *LDAPClient) SyncAllUsers(ctx context.Context) ([]*LDAPUser, error) {
    searchBase := c.config.UserSearchBase
    if searchBase == "" {
        searchBase = c.config.BaseDN
    }

    // Базовый фильтр для пользователей
    filter := "(objectClass=person)"
    
    attributes := []string{"dn"}
    for _, attr := range c.config.AttributeMapping {
        attributes = append(attributes, attr)
    }
    attributes = append(attributes, "memberOf")

    searchRequest := ldap.NewSearchRequest(
        searchBase,
        ldap.ScopeWholeSubtree,
        ldap.NeverDerefAliases,
        0,
        int(c.config.SearchTimeout.Seconds()),
        false,
        filter,
        attributes,
        nil,
    )

    result, err := c.conn.Search(searchRequest)
    if err != nil {
        return nil, fmt.Errorf("LDAP sync search failed: %w", err)
    }

    users := make([]*LDAPUser, 0, len(result.Entries))
    for _, entry := range result.Entries {
        user := c.entryToUser(entry)
        user.Groups = entry.GetAttributeValues("memberOf")
        users = append(users, user)
    }

    return users, nil
}

// entryToUser конвертирует LDAP entry в структуру LDAPUser
func (c *LDAPClient) entryToUser(entry *ldap.Entry) *LDAPUser {
    user := &LDAPUser{
        DN:         entry.DN,
        Attributes: make(map[string]string),
    }

    for field, ldapAttr := range c.config.AttributeMapping {
        value := entry.GetAttributeValue(ldapAttr)
        user.Attributes[field] = value
        
        switch field {
        case "username":
            user.Username = value
        case "email":
            user.Email = value
        case "first_name":
            user.FirstName = value
        case "last_name":
            user.LastName = value
        case "display_name":
            user.DisplayName = value
        case "employee_id":
            user.EmployeeID = value
        case "department":
            user.Department = value
        }
    }

    return user
}

// TestConnection проверяет подключение к LDAP
func (c *LDAPClient) TestConnection() error {
    if err := c.Connect(); err != nil {
        return err
    }
    defer c.Close()
    
    // Пробуем выполнить простой поиск
    searchRequest := ldap.NewSearchRequest(
        c.config.BaseDN,
        ldap.ScopeBaseObject,
        ldap.NeverDerefAliases,
        1,
        10,
        false,
        "(objectClass=*)",
        []string{"dn"},
        nil,
    )
    
    _, err := c.conn.Search(searchRequest)
    return err
}

func (c *LDAPClient) createConnection() (*ldap.Conn, error) {
    address := fmt.Sprintf("%s:%d", c.config.Host, c.config.Port)
    
    if c.config.UseSSL {
        return ldap.DialTLS("tcp", address, &tls.Config{
            ServerName: c.config.Host,
            MinVersion: tls.VersionTLS12,
        })
    }
    
    conn, err := ldap.Dial("tcp", address)
    if err != nil {
        return nil, err
    }
    
    if c.config.UseStartTLS {
        if err := conn.StartTLS(&tls.Config{ServerName: c.config.Host}); err != nil {
            conn.Close()
            return nil, err
        }
    }
    
    return conn, nil
}

// LDAPUser представляет пользователя из LDAP
type LDAPUser struct {
    DN          string
    Username    string
    Email       string
    FirstName   string
    LastName    string
    DisplayName string
    EmployeeID  string
    Department  string
    Groups      []string
    Attributes  map[string]string
}
```

```go
// internal/auth/ldap/service.go
package ldap

import (
    "context"
    "fmt"
    "strings"
    "time"

    "github.com/google/uuid"
)

// LDAPService управляет LDAP-аутентификацией и синхронизацией
type LDAPService struct {
    repo       LDAPRepository
    userRepo   UserRepository
    clients    map[uuid.UUID]*LDAPClient
}

func NewLDAPService(repo LDAPRepository, userRepo UserRepository) *LDAPService {
    return &LDAPService{
        repo:     repo,
        userRepo: userRepo,
        clients:  make(map[uuid.UUID]*LDAPClient),
    }
}

// Authenticate аутентифицирует пользователя через LDAP
func (s *LDAPService) Authenticate(ctx context.Context, tenantID uuid.UUID, username, password string) (*User, error) {
    // 1. Получаем активную LDAP-конфигурацию для tenant
    config, err := s.repo.GetActiveConfig(ctx, tenantID)
    if err != nil {
        return nil, fmt.Errorf("LDAP not configured: %w", err)
    }

    // 2. Получаем или создаем клиент
    client, err := s.getClient(config)
    if err != nil {
        return nil, err
    }

    // 3. Аутентифицируем в LDAP
    ldapUser, err := client.Authenticate(ctx, username, password)
    if err != nil {
        return nil, err
    }

    // 4. Находим или создаем пользователя в нашей системе
    user, err := s.findOrCreateUser(ctx, tenantID, config, ldapUser)
    if err != nil {
        return nil, err
    }

    // 5. Обновляем маппинг и время последнего входа
    s.repo.UpdateLastLogin(ctx, config.ID, ldapUser.DN)

    return user, nil
}

// findOrCreateUser находит или создает пользователя на основе LDAP данных
func (s *LDAPService) findOrCreateUser(
    ctx context.Context,
    tenantID uuid.UUID,
    config *LDAPConfiguration,
    ldapUser *LDAPUser,
) (*User, error) {
    // 1. Ищем по LDAP маппингу
    mapping, err := s.repo.FindMapping(ctx, config.ID, ldapUser.DN)
    if err == nil && mapping != nil {
        user, err := s.userRepo.FindByID(ctx, mapping.UserID)
        if err != nil {
            return nil, err
        }
        
        // Обновляем профиль если настроено
        if config.AutoUpdateProfile {
            s.updateUserFromLDAP(ctx, user, ldapUser, config)
        }
        
        // Обновляем роли если настроено
        if config.SyncGroups {
            s.syncUserRoles(ctx, user, ldapUser.Groups, config.RoleMapping)
        }
        
        return user, nil
    }

    // 2. Ищем по email
    if ldapUser.Email != "" {
        user, err := s.userRepo.FindByEmail(ctx, tenantID, ldapUser.Email)
        if err == nil && user != nil {
            // Связываем существующего пользователя с LDAP
            s.createMapping(ctx, config.ID, user.ID, ldapUser)
            return user, nil
        }
    }

    // 3. Создаем нового пользователя (JIT provisioning)
    if !config.AutoCreateUsers {
        return nil, fmt.Errorf("user not found and auto-creation disabled")
    }

    user, err := s.createUserFromLDAP(ctx, tenantID, config, ldapUser)
    if err != nil {
        return nil, err
    }

    // 4. Создаем маппинг
    s.createMapping(ctx, config.ID, user.ID, ldapUser)

    return user, nil
}

// createUserFromLDAP создает пользователя из LDAP данных
func (s *LDAPService) createUserFromLDAP(
    ctx context.Context,
    tenantID uuid.UUID,
    config *LDAPConfiguration,
    ldapUser *LDAPUser,
) (*User, error) {
    // Определяем роль на основе групп
    role := s.determineRole(ldapUser.Groups, config.RoleMapping)

    user := &User{
        ID:        uuid.New(),
        TenantID:  tenantID,
        Email:     ldapUser.Email,
        Username:  ldapUser.Username,
        FirstName: ldapUser.FirstName,
        LastName:  ldapUser.LastName,
        Role:      role,
        IsActive:  true,
        CreatedAt: time.Now(),
        Source:    "ldap",
    }

    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}

// determineRole определяет роль пользователя на основе групп LDAP
func (s *LDAPService) determineRole(groups []string, roleMapping map[string]string) string {
    // Приоритет: admin > instructor > student
    rolePriority := map[string]int{
        "admin":      3,
        "instructor": 2,
        "student":    1,
    }
    
    highestRole := "student"
    highestPriority := 0

    for _, groupDN := range groups {
        groupDNLower := strings.ToLower(groupDN)
        for mappedGroup, role := range roleMapping {
            if strings.ToLower(mappedGroup) == groupDNLower || 
               strings.Contains(groupDNLower, strings.ToLower(mappedGroup)) {
                if priority, ok := rolePriority[role]; ok && priority > highestPriority {
                    highestRole = role
                    highestPriority = priority
                }
            }
        }
    }

    return highestRole
}

// SyncUsers выполняет полную синхронизацию пользователей
func (s *LDAPService) SyncUsers(ctx context.Context, configID uuid.UUID) (*SyncResult, error) {
    config, err := s.repo.GetConfigByID(ctx, configID)
    if err != nil {
        return nil, err
    }

    client, err := s.getClient(config)
    if err != nil {
        return nil, err
    }

    // Создаем запись о синхронизации
    syncLog := &SyncLog{
        ID:           uuid.New(),
        ConfigID:     configID,
        SyncType:     "full",
        StartedAt:    time.Now(),
        Status:       "running",
    }
    s.repo.CreateSyncLog(ctx, syncLog)

    // Получаем всех пользователей из LDAP
    ldapUsers, err := client.SyncAllUsers(ctx)
    if err != nil {
        syncLog.Status = "failed"
        syncLog.ErrorDetails = map[string]interface{}{"error": err.Error()}
        s.repo.UpdateSyncLog(ctx, syncLog)
        return nil, err
    }

    result := &SyncResult{
        UsersFound: len(ldapUsers),
    }

    // Обрабатываем каждого пользователя
    for _, ldapUser := range ldapUsers {
        if ldapUser.Email == "" {
            continue // Пропускаем пользователей без email
        }

        user, created, err := s.syncUser(ctx, config, ldapUser)
        if err != nil {
            result.Errors = append(result.Errors, err.Error())
            continue
        }

        if created {
            result.UsersCreated++
        } else {
            result.UsersUpdated++
        }
        _ = user
    }

    // Обновляем лог
    syncLog.CompletedAt = time.Now()
    syncLog.Status = "success"
    syncLog.UsersFound = result.UsersFound
    syncLog.UsersCreated = result.UsersCreated
    syncLog.UsersUpdated = result.UsersUpdated
    syncLog.ErrorsCount = len(result.Errors)
    s.repo.UpdateSyncLog(ctx, syncLog)

    return result, nil
}

// TestConnection тестирует подключение к LDAP
func (s *LDAPService) TestConnection(ctx context.Context, config *LDAPConfiguration) error {
    client := NewLDAPClient(&LDAPConfig{
        Host:             config.Host,
        Port:             config.Port,
        UseSSL:           config.UseSSL,
        UseStartTLS:      config.UseStartTLS,
        BindDN:           config.BindDN,
        BindPassword:     config.BindPassword,
        BaseDN:           config.BaseDN,
        ConnectionTimeout: 10 * time.Second,
        SearchTimeout:    30 * time.Second,
    })
    
    return client.TestConnection()
}

func (s *LDAPService) getClient(config *LDAPConfiguration) (*LDAPClient, error) {
    if client, ok := s.clients[config.ID]; ok {
        return client, nil
    }

    client := NewLDAPClient(&LDAPConfig{
        Host:             config.Host,
        Port:             config.Port,
        UseSSL:           config.UseSSL,
        UseStartTLS:      config.UseStartTLS,
        BindDN:           config.BindDN,
        BindPassword:     config.BindPassword,
        BaseDN:           config.BaseDN,
        UserSearchBase:   config.UserSearchBase,
        UserSearchFilter: config.UserSearchFilter,
        GroupSearchBase:  config.GroupSearchBase,
        AttributeMapping: config.AttributeMapping,
        RoleMapping:      config.RoleMapping,
        ConnectionTimeout: time.Duration(config.ConnectionTimeoutSeconds) * time.Second,
        SearchTimeout:    time.Duration(config.SearchTimeoutSeconds) * time.Second,
    })

    if err := client.Connect(); err != nil {
        return nil, err
    }

    s.clients[config.ID] = client
    return client, nil
}

// SyncResult результат синхронизации
type SyncResult struct {
    UsersFound    int
    UsersCreated  int
    UsersUpdated  int
    UsersDisabled int
    Errors        []string
}
```

#### API Endpoints:

```go
// internal/handlers/ldap_handler.go
package handlers

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type LDAPHandler struct {
    ldapService *ldap.LDAPService
}

// POST /api/v1/auth/ldap/login
// Аутентификация через LDAP
func (h *LDAPHandler) Login(c *gin.Context) {
    var req struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }
    
    tenantID := c.GetString("tenant_id")
    
    user, err := h.ldapService.Authenticate(
        c.Request.Context(),
        uuid.MustParse(tenantID),
        req.Username,
        req.Password,
    )
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication failed"})
        return
    }
    
    // Создаем сессию
    token, err := h.sessionService.CreateSession(c.Request.Context(), user)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "session creation failed"})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "user": user,
    })
}

// GET /api/v1/admin/ldap/configurations
// Список LDAP-конфигураций для tenant
func (h *LDAPHandler) ListConfigurations(c *gin.Context) {
    tenantID := c.GetString("tenant_id")
    
    configs, err := h.ldapService.ListConfigurations(c.Request.Context(), uuid.MustParse(tenantID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    // Не возвращаем пароли
    for _, config := range configs {
        config.BindPassword = "********"
    }
    
    c.JSON(http.StatusOK, configs)
}

// POST /api/v1/admin/ldap/configurations
// Создание LDAP-конфигурации
func (h *LDAPHandler) CreateConfiguration(c *gin.Context) {
    var config ldap.LDAPConfiguration
    if err := c.ShouldBindJSON(&config); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    tenantID := c.GetString("tenant_id")
    config.TenantID = uuid.MustParse(tenantID)
    
    if err := h.ldapService.CreateConfiguration(c.Request.Context(), &config); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    config.BindPassword = "********"
    c.JSON(http.StatusCreated, config)
}

// POST /api/v1/admin/ldap/configurations/:id/test
// Тестирование подключения
func (h *LDAPHandler) TestConnection(c *gin.Context) {
    configID := c.Param("id")
    
    config, err := h.ldapService.GetConfiguration(c.Request.Context(), uuid.MustParse(configID))
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "configuration not found"})
        return
    }
    
    if err := h.ldapService.TestConnection(c.Request.Context(), config); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "error":   err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "Connection successful",
    })
}

// POST /api/v1/admin/ldap/configurations/:id/sync
// Запуск синхронизации пользователей
func (h *LDAPHandler) SyncUsers(c *gin.Context) {
    configID := c.Param("id")
    
    result, err := h.ldapService.SyncUsers(c.Request.Context(), uuid.MustParse(configID))
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, result)
}
```

---

### 2.6 Сложность освоения и реализации

#### Оценка сложности:

| Аспект | Сложность | Комментарий |
|--------|-----------|-------------|
| **Базовая аутентификация** | 🟢 Низкая | Bind + Search — простые операции |
| **Active Directory** | 🟡 Средняя | Специфичные атрибуты (sAMAccountName, memberOf) |
| **OpenLDAP** | 🟡 Средняя | Другая схема групп |
| **Синхронизация групп** | 🟡 Средняя | Nested groups в AD — сложно |
| **SSL/TLS** | 🟢 Низкая | Стандартная настройка |
| **Connection Pooling** | 🟡 Средняя | Важно для производительности |

#### Временные оценки:

```
Реализация базовой LDAP-аутентификации:
├── Изучение протокола LDAP: 1-2 дня
├── Базовая реализация (bind, search): 2-3 дня
├── Тестирование с Active Directory: 2-3 дня
├── Тестирование с OpenLDAP: 1-2 дня
└── Итого: 1 неделя

Реализация синхронизации:
├── Полная синхронизация: 2-3 дня
├── Инкрементальная синхронизация: 2-3 дня
├── Обработка групп: 2-3 дня
└── Итого: 1 неделя

Admin UI:
├── Формы конфигурации: 2-3 дня
├── Test connection: 1 день
├── Sync UI: 1-2 дня
└── Итого: 4-5 дней

Общее время: 2-3 недели (один разработчик)
```

#### Типичные проблемы и решения:

| Проблема | Причина | Решение |
|----------|---------|---------|
| "Connection refused" | Firewall блокирует порт | Открыть 389 (LDAP) или 636 (LDAPS) |
| "Invalid credentials" при bind | Неверный формат Bind DN | Использовать полный DN или UPN |
| Пользователь не найден | Неверный search filter | Для AD: `(sAMAccountName=%s)` |
| Группы не синхронизируются | Nested groups в AD | Использовать recursive search |
| Медленная аутентификация | Нет connection pool | Реализовать пул соединений |
| Timeout при поиске | Слишком широкий scope | Использовать более узкий search base |

#### Различия Active Directory vs OpenLDAP:

| Аспект | Active Directory | OpenLDAP |
|--------|------------------|----------|
| **Username attr** | `sAMAccountName` | `uid` |
| **User filter** | `(objectClass=user)` | `(objectClass=inetOrgPerson)` |
| **Group membership** | `memberOf` атрибут у пользователя | `member` атрибут у группы |
| **Nested groups** | Поддерживается нативно | Требует overlay |
| **Password policy** | Встроено | Отдельный модуль |
| **SSL** | Порт 636 | Порт 636 или STARTTLS |

---

### 2.7 Источники для дальнейшего изучения

#### Официальные спецификации:

| Ресурс | Ссылка | Описание |
|--------|--------|----------|
| **RFC 4511** | [tools.ietf.org/html/rfc4511](https://tools.ietf.org/html/rfc4511) | LDAP Protocol |
| **RFC 4512** | [tools.ietf.org/html/rfc4512](https://tools.ietf.org/html/rfc4512) | Directory Information Models |
| **RFC 4513** | [tools.ietf.org/html/rfc4513](https://tools.ietf.org/html/rfc4513) | LDAP Authentication Methods |
| **RFC 4516** | [tools.ietf.org/html/rfc4516](https://tools.ietf.org/html/rfc4516) | LDAP URL Format |

#### Библиотеки (Go):

| Библиотека | Ссылка | Описание |
|------------|--------|----------|
| **go-ldap/ldap** | [github.com/go-ldap/ldap](https://github.com/go-ldap/ldap) | Основная LDAP библиотека для Go |
| **go-asn1-ber** | [github.com/go-asn1-ber/asn1-ber](https://github.com/go-asn1-ber/asn1-ber) | ASN.1 BER encoding |

#### Документация вендоров:

| Ресурс | Ссылка | Описание |
|--------|--------|----------|
| **Microsoft AD** | [docs.microsoft.com/en-us/windows-server/identity/ad-ds](https://docs.microsoft.com/en-us/windows-server/identity/ad-ds/) | Active Directory Documentation |
| **OpenLDAP** | [openldap.org/doc](https://www.openldap.org/doc/) | OpenLDAP Admin Guide |
| **FreeIPA** | [freeipa.org/page/Documentation](https://www.freeipa.org/page/Documentation) | FreeIPA (Red Hat IdM) |

#### Инструменты для тестирования:

| Инструмент | Назначение | Ссылка/Команда |
|------------|------------|----------------|
| **ldapsearch** | CLI для поиска в LDAP | Встроен в OpenLDAP |
| **Apache Directory Studio** | GUI клиент для LDAP | [directory.apache.org/studio](https://directory.apache.org/studio/) |
| **LDAP Admin** | Windows GUI клиент | [ldapadmin.org](http://www.ldapadmin.org/) |
| **Docker OpenLDAP** | Тестовый LDAP-сервер | `docker run osixia/openldap` |

#### Примеры ldapsearch:

```bash
# Поиск пользователя в Active Directory
ldapsearch -H ldap://dc.university.edu -D "cn=admin,dc=university,dc=edu" \
  -w password -b "dc=university,dc=edu" "(sAMAccountName=jdoe)"

# Поиск всех пользователей
ldapsearch -H ldap://ldap.university.edu -D "cn=admin,dc=university,dc=edu" \
  -w password -b "ou=Users,dc=university,dc=edu" "(objectClass=person)"

# Поиск групп пользователя
ldapsearch -H ldap://dc.university.edu -D "cn=admin,dc=university,dc=edu" \
  -w password -b "dc=university,dc=edu" "(member=cn=John Doe,ou=Users,dc=university,dc=edu)"
```

---

### 2.8 Чек-лист реализации

```
Phase 1: Базовая аутентификация (Day 1-3)
□ Установить go-ldap/ldap библиотеку
□ Реализовать Connect/Bind
□ Реализовать SearchUser
□ Реализовать Authenticate (bind as user)
□ Unit тесты с mock LDAP

Phase 2: Active Directory (Day 4-5)
□ Настроить тестовый AD (или использовать клиентский)
□ Протестировать с sAMAccountName
□ Реализовать получение memberOf
□ Обработать nested groups
□ Протестировать LDAPS (порт 636)

Phase 3: Синхронизация (Day 6-8)
□ Реализовать SyncAllUsers
□ Реализовать маппинг групп → роли
□ Реализовать JIT provisioning
□ Реализовать scheduled sync (cron)
□ Логирование синхронизации

Phase 4: Admin UI (Day 9-11)
□ Форма добавления LDAP-конфигурации
□ Test Connection кнопка
□ Sync Now кнопка
□ Просмотр логов синхронизации
□ Маппинг атрибутов UI
□ Маппинг групп → роли UI

Phase 5: Тестирование и документация (Day 12-14)
□ Интеграционные тесты
□ Тестирование с реальным AD клиента
□ Тестирование с OpenLDAP
□ Документация для администраторов
□ Troubleshooting guide
```

---

## 3. WebSocket/Real-time Communication

### 3.1 Определение

**WebSocket** — это протокол связи, обеспечивающий полнодуплексный канал связи через одно TCP-соединение. В отличие от HTTP, где клиент должен инициировать каждый запрос, WebSocket позволяет серверу отправлять данные клиенту в любой момент без предварительного запроса.

**Real-time Communication** — это мгновенный обмен данными между пользователями или между сервером и клиентом с минимальной задержкой (обычно <100ms). Это основа для чатов, уведомлений, совместной работы и интерактивных функций.

#### Сравнение подходов к обновлению данных:

| Подход | Описание | Задержка | Нагрузка на сервер |
|--------|----------|----------|-------------------|
| **Polling** | Клиент периодически запрашивает сервер | 1-30 сек | Высокая (много запросов) |
| **Long Polling** | Сервер держит соединение до появления данных | 0-30 сек | Средняя |
| **Server-Sent Events (SSE)** | Односторонний поток от сервера к клиенту | <100ms | Низкая |
| **WebSocket** | Двусторонний канал в реальном времени | <50ms | Низкая |

#### Ключевые концепции:

| Термин | Описание | Пример |
|--------|----------|--------|
| **Connection** | Постоянное TCP-соединение между клиентом и сервером | ws://app.com/ws |
| **Channel/Room** | Логическая группировка соединений | `course:123`, `chat:456` |
| **Subscription** | Подписка клиента на канал | Студент подписан на курс |
| **Broadcast** | Отправка сообщения всем в канале | Новое сообщение в чате |
| **Presence** | Информация об онлайн-статусе пользователей | "John is online" |
| **Pub/Sub** | Паттерн издатель-подписчик | Redis Pub/Sub |

#### Архитектура WebSocket соединения:

```
┌─────────────────────────────────────────────────────────────────┐
│                    WebSocket Connection Flow                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────┐                              ┌──────────────────┐ │
│  │  Client  │                              │     Server       │ │
│  │ (Browser)│                              │   (Centrifugo)   │ │
│  └────┬─────┘                              └────────┬─────────┘ │
│       │                                             │           │
│       │  1. HTTP Upgrade Request                    │           │
│       │  GET /ws HTTP/1.1                           │           │
│       │  Upgrade: websocket                         │           │
│       │  Connection: Upgrade                        │           │
│       │────────────────────────────────────────────>│           │
│       │                                             │           │
│       │  2. HTTP 101 Switching Protocols            │           │
│       │<────────────────────────────────────────────│           │
│       │                                             │           │
│       │  ═══════ WebSocket Connection Open ═══════  │           │
│       │                                             │           │
│       │  3. Subscribe to channel "course:123"       │           │
│       │────────────────────────────────────────────>│           │
│       │                                             │           │
│       │  4. Subscription confirmed                  │           │
│       │<────────────────────────────────────────────│           │
│       │                                             │           │
│       │  5. Server pushes new message               │           │
│       │<────────────────────────────────────────────│           │
│       │                                             │           │
│       │  6. Client sends message                    │           │
│       │────────────────────────────────────────────>│           │
│       │                                             │           │
│       │  7. Server broadcasts to all subscribers    │           │
│       │<────────────────────────────────────────────│           │
│       │                                             │           │
└─────────────────────────────────────────────────────────────────┘
```

#### Типичная архитектура для масштабирования:

```
┌─────────────────────────────────────────────────────────────────┐
│                     Scalable Real-time Architecture              │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│    Clients                Load Balancer              Backend     │
│                                                                  │
│  ┌────────┐          ┌─────────────────┐                        │
│  │Browser1│─────────>│                 │                        │
│  └────────┘          │                 │      ┌──────────────┐  │
│  ┌────────┐          │   Nginx/HAProxy │      │  API Server  │  │
│  │Browser2│─────────>│   (Sticky)      │<────>│   (Go/Gin)   │  │
│  └────────┘          │                 │      └──────┬───────┘  │
│  ┌────────┐          │                 │             │          │
│  │Browser3│─────────>│                 │             │          │
│  └────────┘          └────────┬────────┘             │          │
│                               │                      │          │
│         ┌─────────────────────┼──────────────────────┘          │
│         │                     │                                  │
│         ▼                     ▼                                  │
│  ┌─────────────┐       ┌─────────────┐       ┌─────────────┐   │
│  │ Centrifugo  │       │ Centrifugo  │       │ Centrifugo  │   │
│  │   Node 1    │       │   Node 2    │       │   Node 3    │   │
│  └──────┬──────┘       └──────┬──────┘       └──────┬──────┘   │
│         │                     │                     │           │
│         └─────────────────────┼─────────────────────┘           │
│                               │                                  │
│                        ┌──────▼──────┐                          │
│                        │    Redis    │                          │
│                        │  (Pub/Sub)  │                          │
│                        └─────────────┘                          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

---

### 3.2 Почему важно реализовать

#### Бизнес-причины:

| Причина | Описание | Влияние на бизнес |
|---------|----------|-------------------|
| **Modern UX Expectation** | Пользователи ожидают мгновенной реакции | Без real-time = устаревший продукт |
| **Engagement** | Мгновенные уведомления увеличивают вовлеченность | +40% DAU с push/real-time |
| **Competitive Parity** | Canvas, Google Classroom имеют real-time | Обязательно для конкуренции |
| **Collaboration** | Совместная работа требует синхронизации | Невозможно без WebSocket |

#### Технические причины:

1. **Эффективность** — один WebSocket заменяет сотни HTTP-запросов
2. **Низкая задержка** — сообщения доставляются за <50ms вместо 1-30 секунд
3. **Масштабируемость** — меньше нагрузка на сервер при большем количестве пользователей
4. **Bidirectional** — сервер может инициировать отправку данных
5. **Battery-friendly** — мобильные устройства экономят батарею без постоянного polling

#### Статистика и факты:

```
📊 Влияние real-time на метрики:
• Slack: 90% сообщений доставляются за <100ms
• Discord: поддерживает 10M+ одновременных подключений
• Google Docs: real-time collaboration используют 1B+ пользователей
• Исследования показывают:
  - +35% времени в приложении при наличии real-time чата
  - +50% retention при push-уведомлениях
  - -70% нагрузки на сервер при переходе с polling на WebSocket
```

#### Сравнение: Polling vs WebSocket

```
Сценарий: 1000 пользователей, обновление каждые 5 секунд

Polling:
├── Запросов в минуту: 1000 × 12 = 12,000
├── Средняя задержка: 2.5 секунды
├── Нагрузка на API: Высокая
└── Батарея мобильного: Быстро разряжается

WebSocket:
├── Запросов в минуту: ~100 (только реальные события)
├── Средняя задержка: <50ms
├── Нагрузка на API: Минимальная
└── Батарея мобильного: Экономится
```

---

### 3.3 Что дает конечному пользователю

#### Для студентов:

| Функция | Без Real-time | С Real-time |
|---------|---------------|-------------|
| **Сообщения в чате** | Обновление через 5-30 сек | Мгновенно (<100ms) |
| **Уведомления** | Видит при обновлении страницы | Всплывает сразу |
| **Статус "печатает"** | Нет | Видит, что собеседник печатает |
| **Онлайн-статус** | Нет | Видит, кто онлайн |
| **Совместная работа** | Невозможна | Редактирование в реальном времени |
| **Оценки** | Узнает при входе в журнал | Push-уведомление сразу |

#### Для преподавателей:

| Функция | Без Real-time | С Real-time |
|---------|---------------|-------------|
| **Вопросы студентов** | Проверяет вручную | Уведомление сразу |
| **Сдача работ** | Узнает при входе | Мгновенное уведомление |
| **Live Q&A** | Невозможно | Интерактивные сессии |
| **Присутствие** | Ручная перекличка | Автоматический трекинг |
| **Polls в лекции** | Сторонние сервисы | Встроенные live polls |

#### Конкретные сценарии использования:

```
Сценарий 1: Чат курса
├─ Без WebSocket: Студент отправляет сообщение → ждет 5 сек → видит ответ
│                 UX как в email 2000-х годов
└─ С WebSocket:   Сообщение появляется мгновенно у всех
                  Видно, кто печатает
                  Индикатор "прочитано" ✓

Сценарий 2: Уведомление об оценке
├─ Без WebSocket: Преподаватель ставит оценку → Студент узнает через 1-24 часа
│                 (когда зайдет в журнал)
└─ С WebSocket:   Оценка появляется → Студент видит уведомление через 100ms
                  Push на телефон, если в фоне ✓

Сценарий 3: Live лекция с Q&A
├─ Без WebSocket: Невозможно реализовать интерактивно
└─ С WebSocket:   Студенты задают вопросы в реальном времени
                  Преподаватель видит поток вопросов
                  Upvote популярных вопросов
                  Live polls с мгновенными результатами ✓

Сценарий 4: Совместное редактирование документа
├─ Без WebSocket: Конфликты при сохранении, потеря данных
└─ С WebSocket:   Google Docs-like experience
                  Курсоры других пользователей видны
                  Изменения синхронизируются в реальном времени ✓
```

#### Engagement метрики:

```
Ожидаемые улучшения после внедрения real-time:

📈 Вовлеченность:
• Сообщений в чате: +200-400%
• Время в приложении: +35%
• DAU/MAU ratio: +15-25%

📱 Удержание:
• 7-day retention: +20%
• Push notification open rate: 40-60%
• Возврат после уведомления: 65%

⚡ Производительность:
• Perceived latency: -80%
• Server load: -60%
• User satisfaction: +40%
```

---

### 3.4 Какие пользователи предпочитают эту функцию

#### Целевые сегменты:

| Сегмент | Важность Real-time | Ключевые use cases |
|---------|-------------------|-------------------|
| **Студенты Gen Z** | Критическая | Чат, уведомления, collaboration |
| **Онлайн-курсы** | Критическая | Live сессии, Q&A, взаимодействие |
| **Корпоративное обучение** | Высокая | Уведомления о дедлайнах, чат с ментором |
| **K-12** | Высокая | Родительские уведомления, безопасность |
| **Blended learning** | Средняя | Синхронизация онлайн/офлайн |

#### Ожидания по возрасту:

```
Поколение Z (1997-2012):
├── Выросли с Instagram, Snapchat, TikTok
├── Ожидают мгновенную реакцию (<1 сек)
├── "Typing indicator" — must have
└── Без real-time = "сломанное" приложение

Миллениалы (1981-1996):
├── Привыкли к email, но адаптировались к Slack
├── Ценят уведомления о важном
└── Готовы ждать 5-10 сек для некритичных функций

Поколение X и старше:
├── Менее требовательны к скорости
├── Ценят уведомления об оценках/дедлайнах
└── Real-time чат — nice to have
```

#### Типичные вопросы при продаже:

```
Вопросы от современных клиентов:
1. "Есть ли real-time чат для студентов?"
2. "Могут ли преподаватели проводить live Q&A?"
3. "Приходят ли push-уведомления об оценках?"
4. "Поддерживается ли совместная работа над документами?"
5. "Какая задержка при отправке сообщений?"

Без WebSocket ответ на первый вопрос = "Нет real-time" = потеря клиента
```

---

### 3.5 Как интегрировать в наше приложение

#### Выбор технологии:

| Опция | Плюсы | Минусы | Рекомендация |
|-------|-------|--------|--------------|
| **Centrifugo** | Battle-tested, масштабируемый, Go-native | Отдельный сервис | ✅ Рекомендуется |
| **gorilla/websocket** | Полный контроль, нет зависимостей | Много кода, сложное масштабирование | Для простых случаев |
| **Socket.io** | Популярный, много клиентов | Node.js, не Go | Не подходит |
| **Pusher/Ably** | SaaS, нет инфраструктуры | Дорого, vendor lock-in | Для MVP |

**Рекомендация: Centrifugo** — production-ready решение, используется крупными компаниями, легко интегрируется с Go.

#### Архитектура с Centrifugo:

```
┌─────────────────────────────────────────────────────────────────┐
│                  Real-time Architecture with Centrifugo          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────┐                         ┌──────────────────┐  │
│  │   Frontend   │                         │   Go Backend     │  │
│  │   (React)    │                         │   (Gin API)      │  │
│  └──────┬───────┘                         └────────┬─────────┘  │
│         │                                          │             │
│         │ WebSocket                                │ HTTP API    │
│         │ (subscribe, receive)                     │ (publish)   │
│         │                                          │             │
│         ▼                                          ▼             │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │                       Centrifugo                            │ │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐ │ │
│  │  │  Channels   │  │  Presence   │  │  History/Recovery   │ │ │
│  │  │  & Pub/Sub  │  │  Tracking   │  │  (missed messages)  │ │ │
│  │  └─────────────┘  └─────────────┘  └─────────────────────┘ │ │
│  └────────────────────────────┬───────────────────────────────┘ │
│                               │                                  │
│                               ▼                                  │
│                        ┌─────────────┐                          │
│                        │    Redis    │                          │
│                        │  (Broker)   │                          │
│                        └─────────────┘                          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘

Каналы (Channels):
├── personal:user_123      — Личные уведомления пользователя
├── course:456             — Все события курса
├── chat:789               — Сообщения чата
├── presence:course:456    — Кто онлайн в курсе
└── typing:chat:789        — Индикаторы печати
```

#### Docker Compose конфигурация:

```yaml
# docker-compose.yml
version: '3.8'

services:
  api:
    build: ./backend
    environment:
      - CENTRIFUGO_API_URL=http://centrifugo:8000/api
      - CENTRIFUGO_API_KEY=${CENTRIFUGO_API_KEY}
      - CENTRIFUGO_TOKEN_SECRET=${JWT_SECRET}
    depends_on:
      - centrifugo
      - redis

  centrifugo:
    image: centrifugo/centrifugo:v5
    command: centrifugo -c config.json
    ports:
      - "8000:8000"  # WebSocket & API
    volumes:
      - ./centrifugo/config.json:/centrifugo/config.json:ro
    environment:
      - CENTRIFUGO_TOKEN_HMAC_SECRET_KEY=${JWT_SECRET}
      - CENTRIFUGO_API_KEY=${CENTRIFUGO_API_KEY}
      - CENTRIFUGO_ADMIN=true
      - CENTRIFUGO_ADMIN_PASSWORD=${CENTRIFUGO_ADMIN_PASSWORD}
    depends_on:
      - redis

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis-data:/data

volumes:
  redis-data:
```

#### Конфигурация Centrifugo:

```json
// centrifugo/config.json
{
  "token_hmac_secret_key": "${CENTRIFUGO_TOKEN_HMAC_SECRET_KEY}",
  "api_key": "${CENTRIFUGO_API_KEY}",
  "admin": true,
  "admin_password": "${CENTRIFUGO_ADMIN_PASSWORD}",
  
  "engine": "redis",
  "redis_address": "redis:6379",
  
  "namespaces": [
    {
      "name": "personal",
      "presence": false,
      "join_leave": false,
      "history_size": 100,
      "history_ttl": "720h",
      "recover": true
    },
    {
      "name": "course",
      "presence": true,
      "join_leave": true,
      "history_size": 100,
      "history_ttl": "168h",
      "recover": true
    },
    {
      "name": "chat",
      "presence": true,
      "join_leave": false,
      "history_size": 500,
      "history_ttl": "720h",
      "recover": true
    },
    {
      "name": "typing",
      "presence": false,
      "join_leave": false,
      "history_size": 0,
      "history_ttl": "0"
    }
  ],
  
  "allowed_origins": [
    "http://localhost:3000",
    "https://app.yourplatform.com"
  ],
  
  "client_channel_limit": 128,
  "channel_max_length": 255,
  
  "websocket_compression": true,
  "websocket_compression_min_size": 128
}
```

#### Схема базы данных:

```sql
-- Хранение сообщений чата (основное хранилище, Centrifugo для доставки)
CREATE TABLE chat_messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    channel_id UUID NOT NULL REFERENCES chat_channels(id),
    
    -- Отправитель
    sender_id UUID NOT NULL REFERENCES users(id),
    
    -- Контент
    content TEXT NOT NULL,
    content_type VARCHAR(20) DEFAULT 'text', -- 'text', 'file', 'image', 'system'
    
    -- Для ответов/тредов
    reply_to_id UUID REFERENCES chat_messages(id),
    thread_root_id UUID REFERENCES chat_messages(id),
    
    -- Файлы/медиа
    attachments JSONB DEFAULT '[]',
    
    -- Метаданные
    metadata JSONB DEFAULT '{}',
    
    -- Статус
    is_edited BOOLEAN DEFAULT false,
    edited_at TIMESTAMP,
    is_deleted BOOLEAN DEFAULT false,
    deleted_at TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT NOW(),
    
    -- Индексы для быстрого поиска
    CONSTRAINT fk_channel FOREIGN KEY (channel_id) REFERENCES chat_channels(id)
);

-- Каналы чата
CREATE TABLE chat_channels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    
    -- Тип канала
    channel_type VARCHAR(20) NOT NULL, -- 'direct', 'group', 'course', 'announcement'
    
    -- Для course channels
    course_id UUID REFERENCES courses(id),
    
    -- Название и описание
    name VARCHAR(200),
    description TEXT,
    
    -- Настройки
    is_private BOOLEAN DEFAULT false,
    allow_reactions BOOLEAN DEFAULT true,
    allow_threads BOOLEAN DEFAULT true,
    
    -- Метаданные
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Участники канала
CREATE TABLE chat_channel_members (
    channel_id UUID NOT NULL REFERENCES chat_channels(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    -- Роль в канале
    role VARCHAR(20) DEFAULT 'member', -- 'owner', 'admin', 'member'
    
    -- Настройки уведомлений
    notifications_enabled BOOLEAN DEFAULT true,
    muted_until TIMESTAMP,
    
    -- Прочитано до
    last_read_at TIMESTAMP,
    last_read_message_id UUID,
    
    -- Метаданные
    joined_at TIMESTAMP DEFAULT NOW(),
    
    PRIMARY KEY (channel_id, user_id)
);

-- Реакции на сообщения
CREATE TABLE message_reactions (
    message_id UUID NOT NULL REFERENCES chat_messages(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    emoji VARCHAR(50) NOT NULL, -- '👍', '❤️', '😂', etc.
    created_at TIMESTAMP DEFAULT NOW(),
    
    PRIMARY KEY (message_id, user_id, emoji)
);

-- Статус набора текста (опционально, можно хранить только в памяти)
CREATE TABLE typing_indicators (
    channel_id UUID NOT NULL,
    user_id UUID NOT NULL,
    started_at TIMESTAMP DEFAULT NOW(),
    
    PRIMARY KEY (channel_id, user_id)
);

-- Уведомления для real-time доставки
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    user_id UUID NOT NULL REFERENCES users(id),
    
    -- Тип уведомления
    type VARCHAR(50) NOT NULL, -- 'grade', 'message', 'announcement', 'deadline', 'submission'
    
    -- Контент
    title VARCHAR(500) NOT NULL,
    body TEXT,
    
    -- Ссылка
    action_url VARCHAR(500),
    
    -- Источник
    source_type VARCHAR(50), -- 'course', 'assignment', 'chat', 'system'
    source_id UUID,
    
    -- Метаданные
    metadata JSONB DEFAULT '{}',
    
    -- Статус
    is_read BOOLEAN DEFAULT false,
    read_at TIMESTAMP,
    is_pushed BOOLEAN DEFAULT false,
    pushed_at TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP
);

-- Индексы
CREATE INDEX idx_chat_messages_channel ON chat_messages(channel_id, created_at DESC);
CREATE INDEX idx_chat_messages_sender ON chat_messages(sender_id);
CREATE INDEX idx_channel_members_user ON chat_channel_members(user_id);
CREATE INDEX idx_notifications_user ON notifications(user_id, is_read, created_at DESC);
CREATE INDEX idx_notifications_type ON notifications(type, created_at DESC);
```

#### Реализация на Go (Backend):

```go
// internal/realtime/centrifugo_client.go
package realtime

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/google/uuid"
)

// CentrifugoClient клиент для взаимодействия с Centrifugo
type CentrifugoClient struct {
    apiURL    string
    apiKey    string
    jwtSecret []byte
    client    *http.Client
}

func NewCentrifugoClient(apiURL, apiKey, jwtSecret string) *CentrifugoClient {
    return &CentrifugoClient{
        apiURL:    apiURL,
        apiKey:    apiKey,
        jwtSecret: []byte(jwtSecret),
        client: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

// GenerateConnectionToken генерирует JWT для подключения клиента к Centrifugo
func (c *CentrifugoClient) GenerateConnectionToken(userID uuid.UUID, expireAt time.Time) (string, error) {
    claims := jwt.MapClaims{
        "sub": userID.String(),
        "exp": expireAt.Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(c.jwtSecret)
}

// GenerateSubscriptionToken генерирует JWT для подписки на канал
func (c *CentrifugoClient) GenerateSubscriptionToken(
    userID uuid.UUID, 
    channel string, 
    expireAt time.Time,
) (string, error) {
    claims := jwt.MapClaims{
        "sub":     userID.String(),
        "channel": channel,
        "exp":     expireAt.Unix(),
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(c.jwtSecret)
}

// Publish отправляет сообщение в канал
func (c *CentrifugoClient) Publish(ctx context.Context, channel string, data interface{}) error {
    payload := map[string]interface{}{
        "channel": channel,
        "data":    data,
    }
    
    return c.apiCall(ctx, "publish", payload)
}

// Broadcast отправляет сообщение в несколько каналов
func (c *CentrifugoClient) Broadcast(ctx context.Context, channels []string, data interface{}) error {
    payload := map[string]interface{}{
        "channels": channels,
        "data":     data,
    }
    
    return c.apiCall(ctx, "broadcast", payload)
}

// Presence получает список пользователей онлайн в канале
func (c *CentrifugoClient) Presence(ctx context.Context, channel string) (*PresenceResult, error) {
    payload := map[string]interface{}{
        "channel": channel,
    }
    
    var result PresenceResult
    if err := c.apiCallWithResult(ctx, "presence", payload, &result); err != nil {
        return nil, err
    }
    
    return &result, nil
}

// PresenceStats получает статистику присутствия
func (c *CentrifugoClient) PresenceStats(ctx context.Context, channel string) (*PresenceStats, error) {
    payload := map[string]interface{}{
        "channel": channel,
    }
    
    var result PresenceStats
    if err := c.apiCallWithResult(ctx, "presence_stats", payload, &result); err != nil {
        return nil, err
    }
    
    return &result, nil
}

// History получает историю сообщений канала
func (c *CentrifugoClient) History(ctx context.Context, channel string, limit int) (*HistoryResult, error) {
    payload := map[string]interface{}{
        "channel": channel,
        "limit":   limit,
    }
    
    var result HistoryResult
    if err := c.apiCallWithResult(ctx, "history", payload, &result); err != nil {
        return nil, err
    }
    
    return &result, nil
}

// Disconnect отключает пользователя от Centrifugo
func (c *CentrifugoClient) Disconnect(ctx context.Context, userID string) error {
    payload := map[string]interface{}{
        "user": userID,
    }
    
    return c.apiCall(ctx, "disconnect", payload)
}

func (c *CentrifugoClient) apiCall(ctx context.Context, method string, params interface{}) error {
    _, err := c.apiCallRaw(ctx, method, params)
    return err
}

func (c *CentrifugoClient) apiCallWithResult(ctx context.Context, method string, params interface{}, result interface{}) error {
    respBody, err := c.apiCallRaw(ctx, method, params)
    if err != nil {
        return err
    }
    
    var resp struct {
        Result json.RawMessage `json:"result"`
    }
    
    if err := json.Unmarshal(respBody, &resp); err != nil {
        return err
    }
    
    return json.Unmarshal(resp.Result, result)
}

func (c *CentrifugoClient) apiCallRaw(ctx context.Context, method string, params interface{}) ([]byte, error) {
    body := map[string]interface{}{
        "method": method,
        "params": params,
    }
    
    jsonBody, err := json.Marshal(body)
    if err != nil {
        return nil, err
    }
    
    req, err := http.NewRequestWithContext(ctx, "POST", c.apiURL, bytes.NewReader(jsonBody))
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "apikey "+c.apiKey)
    
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    var buf bytes.Buffer
    if _, err := buf.ReadFrom(resp.Body); err != nil {
        return nil, err
    }
    
    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("centrifugo API error: %s", buf.String())
    }
    
    return buf.Bytes(), nil
}

// Types
type PresenceResult struct {
    Presence map[string]ClientInfo `json:"presence"`
}

type ClientInfo struct {
    User   string `json:"user"`
    Client string `json:"client"`
}

type PresenceStats struct {
    NumClients int `json:"num_clients"`
    NumUsers   int `json:"num_users"`
}

type HistoryResult struct {
    Publications []Publication `json:"publications"`
}

type Publication struct {
    Data   json.RawMessage `json:"data"`
    Offset uint64          `json:"offset"`
}
```

```go
// internal/realtime/service.go
package realtime

import (
    "context"
    "fmt"
    "time"

    "github.com/google/uuid"
)

// RealtimeService сервис для real-time функций
type RealtimeService struct {
    centrifugo  *CentrifugoClient
    chatRepo    ChatRepository
    notifyRepo  NotificationRepository
}

func NewRealtimeService(
    centrifugo *CentrifugoClient,
    chatRepo ChatRepository,
    notifyRepo NotificationRepository,
) *RealtimeService {
    return &RealtimeService{
        centrifugo: centrifugo,
        chatRepo:   chatRepo,
        notifyRepo: notifyRepo,
    }
}

// GetConnectionCredentials возвращает токен для подключения к Centrifugo
func (s *RealtimeService) GetConnectionCredentials(ctx context.Context, userID uuid.UUID) (*ConnectionCredentials, error) {
    expireAt := time.Now().Add(24 * time.Hour)
    
    token, err := s.centrifugo.GenerateConnectionToken(userID, expireAt)
    if err != nil {
        return nil, err
    }
    
    return &ConnectionCredentials{
        Token:     token,
        ExpiresAt: expireAt,
    }, nil
}

// SendChatMessage отправляет сообщение в чат
func (s *RealtimeService) SendChatMessage(ctx context.Context, msg *ChatMessage) error {
    // 1. Сохраняем в БД
    if err := s.chatRepo.SaveMessage(ctx, msg); err != nil {
        return err
    }
    
    // 2. Публикуем в Centrifugo
    channel := fmt.Sprintf("chat:%s", msg.ChannelID)
    
    payload := map[string]interface{}{
        "type":       "message",
        "id":         msg.ID,
        "sender_id":  msg.SenderID,
        "content":    msg.Content,
        "created_at": msg.CreatedAt,
    }
    
    return s.centrifugo.Publish(ctx, channel, payload)
}

// SendTypingIndicator отправляет индикатор печати
func (s *RealtimeService) SendTypingIndicator(ctx context.Context, channelID, userID uuid.UUID, isTyping bool) error {
    channel := fmt.Sprintf("typing:%s", channelID)
    
    payload := map[string]interface{}{
        "type":      "typing",
        "user_id":   userID,
        "is_typing": isTyping,
    }
    
    return s.centrifugo.Publish(ctx, channel, payload)
}

// SendNotification отправляет уведомление пользователю
func (s *RealtimeService) SendNotification(ctx context.Context, notification *Notification) error {
    // 1. Сохраняем в БД
    if err := s.notifyRepo.Save(ctx, notification); err != nil {
        return err
    }
    
    // 2. Публикуем в персональный канал
    channel := fmt.Sprintf("personal:%s", notification.UserID)
    
    payload := map[string]interface{}{
        "type":       "notification",
        "id":         notification.ID,
        "title":      notification.Title,
        "body":       notification.Body,
        "action_url": notification.ActionURL,
        "created_at": notification.CreatedAt,
    }
    
    return s.centrifugo.Publish(ctx, channel, payload)
}

// BroadcastCourseEvent отправляет событие всем участникам курса
func (s *RealtimeService) BroadcastCourseEvent(ctx context.Context, courseID uuid.UUID, event *CourseEvent) error {
    channel := fmt.Sprintf("course:%s", courseID)
    
    payload := map[string]interface{}{
        "type": event.Type,
        "data": event.Data,
    }
    
    return s.centrifugo.Publish(ctx, channel, payload)
}

// GetOnlineUsers возвращает список пользователей онлайн в канале
func (s *RealtimeService) GetOnlineUsers(ctx context.Context, channelID uuid.UUID) ([]uuid.UUID, error) {
    channel := fmt.Sprintf("presence:course:%s", channelID)
    
    result, err := s.centrifugo.Presence(ctx, channel)
    if err != nil {
        return nil, err
    }
    
    userIDs := make([]uuid.UUID, 0, len(result.Presence))
    for _, client := range result.Presence {
        if uid, err := uuid.Parse(client.User); err == nil {
            userIDs = append(userIDs, uid)
        }
    }
    
    return userIDs, nil
}

// Types
type ConnectionCredentials struct {
    Token     string    `json:"token"`
    ExpiresAt time.Time `json:"expires_at"`
}

type ChatMessage struct {
    ID        uuid.UUID `json:"id"`
    ChannelID uuid.UUID `json:"channel_id"`
    SenderID  uuid.UUID `json:"sender_id"`
    Content   string    `json:"content"`
    CreatedAt time.Time `json:"created_at"`
}

type Notification struct {
    ID        uuid.UUID `json:"id"`
    UserID    uuid.UUID `json:"user_id"`
    Title     string    `json:"title"`
    Body      string    `json:"body"`
    ActionURL string    `json:"action_url"`
    CreatedAt time.Time `json:"created_at"`
}

type CourseEvent struct {
    Type string                 `json:"type"`
    Data map[string]interface{} `json:"data"`
}
```

#### API Handlers:

```go
// internal/handlers/realtime_handler.go
package handlers

import (
    "net/http"
    
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type RealtimeHandler struct {
    realtimeService *realtime.RealtimeService
}

// GET /api/v1/realtime/credentials
// Возвращает токен для подключения к WebSocket
func (h *RealtimeHandler) GetCredentials(c *gin.Context) {
    userID := c.GetString("user_id")
    
    creds, err := h.realtimeService.GetConnectionCredentials(
        c.Request.Context(),
        uuid.MustParse(userID),
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, creds)
}

// POST /api/v1/chat/messages
// Отправка сообщения в чат
func (h *RealtimeHandler) SendMessage(c *gin.Context) {
    var req struct {
        ChannelID string `json:"channel_id" binding:"required"`
        Content   string `json:"content" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    userID := c.GetString("user_id")
    
    msg := &realtime.ChatMessage{
        ID:        uuid.New(),
        ChannelID: uuid.MustParse(req.ChannelID),
        SenderID:  uuid.MustParse(userID),
        Content:   req.Content,
    }
    
    if err := h.realtimeService.SendChatMessage(c.Request.Context(), msg); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusCreated, msg)
}

// POST /api/v1/chat/typing
// Индикатор печати
func (h *RealtimeHandler) SendTyping(c *gin.Context) {
    var req struct {
        ChannelID string `json:"channel_id" binding:"required"`
        IsTyping  bool   `json:"is_typing"`
    }
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    userID := c.GetString("user_id")
    
    if err := h.realtimeService.SendTypingIndicator(
        c.Request.Context(),
        uuid.MustParse(req.ChannelID),
        uuid.MustParse(userID),
        req.IsTyping,
    ); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.Status(http.StatusOK)
}

// GET /api/v1/courses/:id/online
// Получить список онлайн пользователей
func (h *RealtimeHandler) GetOnlineUsers(c *gin.Context) {
    courseID := c.Param("id")
    
    users, err := h.realtimeService.GetOnlineUsers(
        c.Request.Context(),
        uuid.MustParse(courseID),
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, gin.H{"users": users})
}
```

#### Frontend интеграция (React):

```typescript
// frontend/src/lib/centrifuge.ts
import { Centrifuge, Subscription, PublicationContext } from 'centrifuge';

class RealtimeClient {
  private centrifuge: Centrifuge | null = null;
  private subscriptions: Map<string, Subscription> = new Map();
  
  async connect(token: string) {
    this.centrifuge = new Centrifuge('wss://app.yourplatform.com/connection/websocket', {
      token,
    });
    
    this.centrifuge.on('connecting', (ctx) => {
      console.log('Connecting to WebSocket...', ctx);
    });
    
    this.centrifuge.on('connected', (ctx) => {
      console.log('Connected to WebSocket', ctx);
    });
    
    this.centrifuge.on('disconnected', (ctx) => {
      console.log('Disconnected from WebSocket', ctx);
    });
    
    this.centrifuge.connect();
  }
  
  subscribe(channel: string, onMessage: (data: any) => void): Subscription {
    if (!this.centrifuge) {
      throw new Error('Not connected');
    }
    
    const sub = this.centrifuge.newSubscription(channel);
    
    sub.on('publication', (ctx: PublicationContext) => {
      onMessage(ctx.data);
    });
    
    sub.on('subscribing', () => {
      console.log(`Subscribing to ${channel}...`);
    });
    
    sub.on('subscribed', () => {
      console.log(`Subscribed to ${channel}`);
    });
    
    sub.subscribe();
    this.subscriptions.set(channel, sub);
    
    return sub;
  }
  
  unsubscribe(channel: string) {
    const sub = this.subscriptions.get(channel);
    if (sub) {
      sub.unsubscribe();
      this.subscriptions.delete(channel);
    }
  }
  
  disconnect() {
    this.subscriptions.forEach((sub) => sub.unsubscribe());
    this.subscriptions.clear();
    this.centrifuge?.disconnect();
  }
}

export const realtimeClient = new RealtimeClient();
```

```typescript
// frontend/src/hooks/useChat.ts
import { useEffect, useState, useCallback } from 'react';
import { realtimeClient } from '@/lib/centrifuge';

interface ChatMessage {
  id: string;
  sender_id: string;
  content: string;
  created_at: string;
}

interface TypingUser {
  user_id: string;
  is_typing: boolean;
}

export function useChat(channelId: string) {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [typingUsers, setTypingUsers] = useState<string[]>([]);
  
  useEffect(() => {
    // Subscribe to chat messages
    const chatSub = realtimeClient.subscribe(`chat:${channelId}`, (data) => {
      if (data.type === 'message') {
        setMessages((prev) => [...prev, data]);
      }
    });
    
    // Subscribe to typing indicators
    const typingSub = realtimeClient.subscribe(`typing:${channelId}`, (data: TypingUser) => {
      if (data.is_typing) {
        setTypingUsers((prev) => [...new Set([...prev, data.user_id])]);
      } else {
        setTypingUsers((prev) => prev.filter((id) => id !== data.user_id));
      }
    });
    
    return () => {
      realtimeClient.unsubscribe(`chat:${channelId}`);
      realtimeClient.unsubscribe(`typing:${channelId}`);
    };
  }, [channelId]);
  
  const sendMessage = useCallback(async (content: string) => {
    await fetch('/api/v1/chat/messages', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ channel_id: channelId, content }),
    });
  }, [channelId]);
  
  const sendTyping = useCallback(async (isTyping: boolean) => {
    await fetch('/api/v1/chat/typing', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ channel_id: channelId, is_typing: isTyping }),
    });
  }, [channelId]);
  
  return { messages, typingUsers, sendMessage, sendTyping };
}
```

---

### 3.6 Сложность освоения и реализации

#### Оценка сложности:

| Аспект | Сложность | Комментарий |
|--------|-----------|-------------|
| **Centrifugo setup** | 🟢 Низкая | Docker, простая конфигурация |
| **Базовый чат** | 🟢 Низкая | Publish/Subscribe — просто |
| **Presence (кто онлайн)** | 🟡 Средняя | Настройка каналов |
| **Typing indicators** | 🟢 Низкая | Отдельный канал |
| **Масштабирование** | 🟡 Средняя | Redis cluster, load balancing |
| **Offline sync/recovery** | 🟡 Средняя | History, recovery в Centrifugo |
| **Frontend интеграция** | 🟢 Низкая | Centrifuge-js хорошо документирован |

#### Временные оценки:

```
Centrifugo setup и базовая интеграция:
├── Docker setup: 1 день
├── JWT токены для авторизации: 1 день
├── Backend publish API: 2 дня
└── Итого: 4 дня

Чат функционал:
├── База данных (схема, репозиторий): 2 дня
├── Send/receive messages: 2 дня
├── Frontend chat UI: 3 дня
├── Typing indicators: 1 день
└── Итого: 8 дней

Уведомления:
├── Notification service: 2 дня
├── Personal channels: 1 день
├── Frontend notification UI: 2 дня
└── Итого: 5 дней

Presence (кто онлайн):
├── Backend presence API: 1 день
├── Frontend integration: 1 день
└── Итого: 2 дня

Общее время: 3-4 недели (один разработчик)
```

#### Типичные проблемы и решения:

| Проблема | Причина | Решение |
|----------|---------|---------|
| WebSocket не подключается | CORS, proxy | Настроить allowed_origins, nginx proxy |
| Сообщения теряются | Нет recovery | Включить history и recover |
| Медленное подключение | Холодный старт | Connection pooling, keep-alive |
| Дублирование сообщений | Reconnect без offset | Использовать epoch и offset |
| Высокая нагрузка | Много мелких сообщений | Batching, дебаунс typing |

---

### 3.7 Источники для дальнейшего изучения

#### Официальная документация:

| Ресурс | Ссылка | Описание |
|--------|--------|----------|
| **Centrifugo Docs** | [centrifugal.dev](https://centrifugal.dev/) | Полная документация |
| **Centrifuge-js** | [github.com/centrifugal/centrifuge-js](https://github.com/centrifugal/centrifuge-js) | JavaScript клиент |
| **WebSocket RFC 6455** | [tools.ietf.org/html/rfc6455](https://tools.ietf.org/html/rfc6455) | Спецификация протокола |

#### Библиотеки:

| Библиотека | Ссылка | Описание |
|------------|--------|----------|
| **gorilla/websocket** | [github.com/gorilla/websocket](https://github.com/gorilla/websocket) | Go WebSocket library |
| **centrifuge** | [github.com/centrifugal/centrifuge](https://github.com/centrifugal/centrifuge) | Go real-time library |

#### Обучающие ресурсы:

| Ресурс | Формат | Описание |
|--------|--------|----------|
| **Centrifugo Tutorial** | [centrifugal.dev/docs/tutorial](https://centrifugal.dev/docs/tutorial/intro) | Пошаговый туториал |
| **Real-time Web Apps** | Статьи | Паттерны проектирования |

---

### 3.8 Чек-лист реализации

```
Phase 1: Infrastructure (Day 1-3)
□ Docker setup для Centrifugo
□ Redis для Pub/Sub
□ Настройка namespaces
□ JWT token generation
□ Nginx proxy для WebSocket

Phase 2: Backend Integration (Day 4-8)
□ CentrifugoClient реализация
□ Publish API
□ Presence API
□ Database schema для chat
□ Chat repository
□ Chat service
□ API handlers

Phase 3: Notifications (Day 9-11)
□ Notification schema
□ Notification service
□ Personal channels
□ Broadcast по курсу
□ Integration с существующими событиями

Phase 4: Frontend (Day 12-18)
□ Centrifuge-js интеграция
□ Connection management
□ useChat hook
□ useNotifications hook
□ Chat UI component
□ Typing indicators UI
□ Online users UI
□ Notification toast/dropdown

Phase 5: Testing & Polish (Day 19-21)
□ Load testing (1000+ connections)
□ Reconnection handling
□ Offline message recovery
□ Error handling
□ Документация
```

---

## 4. SCORM Support (Sharable Content Object Reference Model)

### 4.1 Определение

**SCORM (Sharable Content Object Reference Model)** — это набор технических стандартов для e-learning, определяющий как создавать, упаковывать и запускать образовательный контент в системах управления обучением (LMS). SCORM обеспечивает совместимость контента между различными LMS-платформами.

**Ключевая идея:** "Write once, run anywhere" — контент, созданный по стандарту SCORM, можно загрузить и запустить в любой SCORM-совместимой LMS без модификаций.

#### Версии SCORM:

| Версия | Год | Особенности | Использование |
|--------|-----|-------------|---------------|
| **SCORM 1.1** | 2001 | Первая версия | Устарела |
| **SCORM 1.2** | 2001 | Стабильная, широко распространена | ~60% контента |
| **SCORM 2004 (1st-4th Ed.)** | 2004-2009 | Sequencing, навигация | ~35% контента |
| **xAPI (Tin Can)** | 2013 | Современный преемник | ~5%, растет |

#### Ключевые компоненты SCORM:

| Компонент | Описание | Пример |
|-----------|----------|--------|
| **SCO (Sharable Content Object)** | Отдельный учебный модуль | Один урок или тема |
| **Asset** | Статический ресурс (изображение, видео) | logo.png, video.mp4 |
| **Manifest (imsmanifest.xml)** | Описание структуры пакета | XML с метаданными |
| **PIF (Package Interchange Format)** | ZIP-архив с контентом | course.zip |
| **RTE (Run-Time Environment)** | API для взаимодействия с LMS | JavaScript API |

#### Структура SCORM-пакета:

```
course.zip
├── imsmanifest.xml          ← Главный файл манифеста
├── adlcp_rootv1p2.xsd       ← Схемы валидации
├── ims_xml.xsd
├── imscp_rootv1p1p2.xsd
├── imsmd_rootv1p2p1.xsd
├── content/
│   ├── module1/
│   │   ├── index.html       ← SCO (запускаемый контент)
│   │   ├── scorm_api.js     ← Wrapper для API
│   │   ├── styles.css
│   │   └── images/
│   ├── module2/
│   │   └── index.html
│   └── shared/
│       └── video.mp4        ← Asset
└── sequencing.xml           ← Для SCORM 2004
```

#### Пример imsmanifest.xml:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<manifest identifier="com.example.course" version="1.0"
          xmlns="http://www.imsproject.org/xsd/imscp_rootv1p1p2"
          xmlns:adlcp="http://www.adlnet.org/xsd/adlcp_rootv1p2">
  
  <metadata>
    <schema>ADL SCORM</schema>
    <schemaversion>1.2</schemaversion>
  </metadata>
  
  <organizations default="org1">
    <organization identifier="org1">
      <title>Введение в программирование</title>
      
      <item identifier="item1" identifierref="res1">
        <title>Модуль 1: Основы</title>
      </item>
      
      <item identifier="item2" identifierref="res2">
        <title>Модуль 2: Переменные</title>
        <adlcp:prerequisites>item1</adlcp:prerequisites>
      </item>
      
      <item identifier="item3" identifierref="res3">
        <title>Финальный тест</title>
        <adlcp:prerequisites>item2</adlcp:prerequisites>
        <adlcp:maxtimeallowed>00:30:00</adlcp:maxtimeallowed>
      </item>
      
    </organization>
  </organizations>
  
  <resources>
    <resource identifier="res1" type="webcontent" 
              adlcp:scormtype="sco" href="content/module1/index.html">
      <file href="content/module1/index.html"/>
      <file href="content/module1/scorm_api.js"/>
      <file href="content/module1/styles.css"/>
    </resource>
    
    <resource identifier="res2" type="webcontent"
              adlcp:scormtype="sco" href="content/module2/index.html">
      <file href="content/module2/index.html"/>
    </resource>
    
    <resource identifier="res3" type="webcontent"
              adlcp:scormtype="sco" href="content/quiz/index.html">
      <file href="content/quiz/index.html"/>
    </resource>
  </resources>
  
</manifest>
```

#### SCORM Run-Time Environment (RTE):

```
┌─────────────────────────────────────────────────────────────────┐
│                    SCORM Runtime Environment                     │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌──────────────────────────────────────────────────────────┐   │
│  │                      LMS (Our Platform)                   │   │
│  │  ┌─────────────────────────────────────────────────────┐ │   │
│  │  │                   SCORM Player                       │ │   │
│  │  │  ┌─────────────┐  ┌────────────────────────────────┐│ │   │
│  │  │  │   iframe    │  │        JavaScript API          ││ │   │
│  │  │  │             │  │  ┌──────────────────────────┐  ││ │   │
│  │  │  │   SCO       │  │  │ window.API (SCORM 1.2)   │  ││ │   │
│  │  │  │  Content    │◄─┼─►│ window.API_1484_11       │  ││ │   │
│  │  │  │  (HTML)     │  │  │      (SCORM 2004)        │  ││ │   │
│  │  │  │             │  │  └──────────────────────────┘  ││ │   │
│  │  │  └─────────────┘  └────────────────────────────────┘│ │   │
│  │  └─────────────────────────────────────────────────────┘ │   │
│  │                            │                              │   │
│  │                            ▼                              │   │
│  │  ┌─────────────────────────────────────────────────────┐ │   │
│  │  │                   Backend API                        │ │   │
│  │  │  POST /api/scorm/initialize                          │ │   │
│  │  │  POST /api/scorm/commit                              │ │   │
│  │  │  POST /api/scorm/terminate                           │ │   │
│  │  └─────────────────────────────────────────────────────┘ │   │
│  └──────────────────────────────────────────────────────────┘   │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

#### SCORM API методы:

**SCORM 1.2 API:**
```javascript
// Инициализация
API.LMSInitialize("")        // Начало сессии
API.LMSFinish("")            // Завершение сессии

// Чтение/запись данных
API.LMSGetValue(element)     // Получить значение
API.LMSSetValue(element, value) // Установить значение
API.LMSCommit("")            // Сохранить данные

// Обработка ошибок
API.LMSGetLastError()        // Код последней ошибки
API.LMSGetErrorString(code)  // Описание ошибки
API.LMSGetDiagnostic(code)   // Детали ошибки
```

**Ключевые Data Model Elements (SCORM 1.2):**

| Element | Описание | Пример значения |
|---------|----------|-----------------|
| `cmi.core.student_id` | ID студента | "user_12345" |
| `cmi.core.student_name` | Имя студента | "Иванов, Иван" |
| `cmi.core.lesson_status` | Статус урока | "completed", "passed", "failed" |
| `cmi.core.score.raw` | Набранные баллы | "85" |
| `cmi.core.score.min` | Минимум баллов | "0" |
| `cmi.core.score.max` | Максимум баллов | "100" |
| `cmi.core.session_time` | Время сессии | "00:15:30" |
| `cmi.core.total_time` | Общее время | "01:45:20" |
| `cmi.suspend_data` | Данные для возобновления | JSON строка |
| `cmi.interactions.n.*` | Ответы на вопросы | Детали взаимодействий |

---

### 4.2 Почему важно реализовать

#### Бизнес-причины:

| Причина | Описание | Влияние на бизнес |
|---------|----------|-------------------|
| **Огромная библиотека контента** | Миллионы готовых SCORM-курсов | Мгновенный доступ к контенту |
| **Индустриальный стандарт** | Все LMS-конкуренты поддерживают | Без SCORM = неконкурентоспособны |
| **Vendor Independence** | Клиенты не хотят vendor lock-in | Возможность экспорта/импорта контента |
| **Corporate Training** | Корпоративное обучение = SCORM | Крупный сегмент рынка |
| **Compliance Training** | Обязательное обучение часто в SCORM | Регуляторные требования |

#### Технические причины:

1. **Переиспользование контента** — один курс работает в любой LMS
2. **Стандартизированный трекинг** — единый формат данных о прогрессе
3. **Богатые возможности** — интерактивность, видео, симуляции
4. **Инвестиции защищены** — контент не устаревает при смене LMS

#### Статистика и факты:

```
📊 Рыночные данные:
• 90%+ корпоративных LMS поддерживают SCORM
• Рынок SCORM-контента: $15+ миллиардов
• 70% корпоративного e-learning в формате SCORM
• Articulate, Adobe Captivate, iSpring — все создают SCORM
• Средняя организация имеет 50-500 SCORM-курсов

💰 Примеры библиотек SCORM-контента:
• LinkedIn Learning: 16,000+ курсов
• Skillsoft: 150,000+ курсов
• OpenSesame: 30,000+ курсов
• Coursera for Business: 5,000+ курсов
```

#### Сравнение: с SCORM vs без SCORM

```
Клиент с существующей библиотекой SCORM:

Без поддержки SCORM:
├── "Нам придется пересоздавать весь контент?"
├── "Это 200 курсов, на которые мы потратили $500K"
├── Решение: Отказ от нашей платформы ❌
└── Потерянный контракт: $50-100K/год

С поддержкой SCORM:
├── "Отлично, мы просто загрузим наши курсы"
├── Миграция за 1 день
├── Решение: Подписание контракта ✅
└── Выигранный контракт: $50-100K/год
```

---

### 4.3 Что дает конечному пользователю

#### Для студентов:

| Функция | Описание | Пользовательский опыт |
|---------|----------|----------------------|
| **Интерактивные курсы** | Видео, симуляции, игры | Engaging обучение |
| **Сохранение прогресса** | Продолжить с того же места | Удобство |
| **Оффлайн прогресс** | SCORM синхронизируется при reconnect | Мобильность |
| **Разнообразный контент** | Разные провайдеры, стили | Выбор |
| **Профессиональные курсы** | LinkedIn Learning, Coursera | Карьерный рост |

#### Для преподавателей/администраторов:

| Функция | Описание | Ценность |
|---------|----------|----------|
| **Готовый контент** | Не нужно создавать с нуля | Экономия времени |
| **Детальная аналитика** | Время, попытки, ответы | Insights |
| **Compliance tracking** | Автоматический трекинг прохождения | Отчетность |
| **Быстрый деплой** | Загрузил ZIP — курс готов | Скорость |
| **Authoring tools** | Articulate, Captivate, iSpring | Создание контента |

#### Конкретные сценарии использования:

```
Сценарий 1: Корпоративный onboarding
├─ HR создает курс в Articulate Storyline
├─ Экспортирует в SCORM 1.2
├─ Загружает в нашу LMS
├─ Новые сотрудники проходят курс
├─ LMS отслеживает: кто прошел, сколько времени, какой балл
└─ HR видит отчет о прохождении ✓

Сценарий 2: Compliance обучение (охрана труда)
├─ Компания покупает готовый SCORM-курс
├─ Загружает в LMS
├─ Назначает всем сотрудникам с дедлайном
├─ LMS отслеживает прохождение
├─ Автоматические напоминания
├─ Сертификат по завершении
└─ Отчет для проверяющих органов ✓

Сценарий 3: Blended learning в университете
├─ Преподаватель находит SCORM-модуль по теме
├─ Встраивает в свой курс
├─ Студенты изучают интерактивный контент
├─ Оценки автоматически попадают в журнал
└─ Преподаватель видит, кто что изучил ✓

Сценарий 4: LinkedIn Learning интеграция
├─ Организация подписана на LinkedIn Learning
├─ Курсы доступны как SCORM через LTI
├─ Сотрудники проходят курсы
├─ Прогресс синхронизируется в нашу LMS
└─ Единый отчет по всему обучению ✓
```

---

### 4.4 Какие пользователи предпочитают эту функцию

#### Целевые сегменты:

| Сегмент | Важность SCORM | Типичный объем контента |
|---------|----------------|------------------------|
| **Корпорации** | Критическая | 100-1000 курсов |
| **Compliance-heavy отрасли** | Критическая | 50-200 обязательных курсов |
| **Университеты** | Высокая | 20-100 курсов |
| **K-12** | Средняя | 10-50 курсов |
| **Стартапы** | Низкая | Мало legacy контента |

#### Отрасли с высоким использованием SCORM:

```
🏥 Здравоохранение:
• HIPAA compliance training
• Medical device training
• Continuing Medical Education (CME)

🏦 Финансы:
• Anti-money laundering (AML)
• Know Your Customer (KYC)
• Regulatory compliance

🏭 Производство:
• Safety training (OSHA)
• Equipment operation
• Quality management

💼 Ритейл:
• Product knowledge
• Customer service
• POS training

✈️ Авиация:
• Pilot training modules
• Safety procedures
• Regulatory compliance
```

#### Типичные вопросы при продаже:

```
Вопросы от корпоративных клиентов:
1. "Поддерживаете ли вы SCORM 1.2 и 2004?"
2. "Можем ли мы загрузить наши существующие курсы?"
3. "Работает ли трекинг completion и score?"
4. "Поддерживается ли suspend/resume?"
5. "Можете ли вы интегрироваться с нашим Articulate/Captivate контентом?"

Без SCORM ответ = "Нет" = потеря 90% корпоративных клиентов
```

---

### 4.5 Как интегрировать в наше приложение

#### Архитектура SCORM Player:

```
┌─────────────────────────────────────────────────────────────────┐
│                    SCORM Integration Architecture                │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │                     Upload & Parse                          │ │
│  │  ┌──────────┐    ┌──────────────┐    ┌─────────────────┐  │ │
│  │  │  Upload  │───►│  Unzip &     │───►│  Parse          │  │ │
│  │  │  ZIP     │    │  Validate    │    │  Manifest       │  │ │
│  │  └──────────┘    └──────────────┘    └────────┬────────┘  │ │
│  │                                               │            │ │
│  │                                    ┌──────────▼─────────┐ │ │
│  │                                    │ Store in S3/MinIO  │ │ │
│  │                                    │ + DB metadata      │ │ │
│  │                                    └────────────────────┘ │ │
│  └────────────────────────────────────────────────────────────┘ │
│                                                                  │
│  ┌────────────────────────────────────────────────────────────┐ │
│  │                     SCORM Player                            │ │
│  │                                                              │ │
│  │  ┌──────────────────────────────────────────────────────┐  │ │
│  │  │                    Browser                            │  │ │
│  │  │  ┌────────────────────────────────────────────────┐  │  │ │
│  │  │  │              Player Container                   │  │  │ │
│  │  │  │  ┌──────────────────┐  ┌────────────────────┐  │  │  │ │
│  │  │  │  │   Navigation     │  │   SCO iframe       │  │  │  │ │
│  │  │  │  │   (TOC, prev,    │  │   ┌────────────┐   │  │  │  │ │
│  │  │  │  │    next)         │  │   │  Content   │   │  │  │  │ │
│  │  │  │  └──────────────────┘  │   │  (HTML)    │   │  │  │  │ │
│  │  │  │                        │   └────────────┘   │  │  │  │ │
│  │  │  │  ┌──────────────────────────────────────┐   │  │  │  │ │
│  │  │  │  │        SCORM API Adapter             │   │  │  │  │ │
│  │  │  │  │  window.API / window.API_1484_11     │   │  │  │  │ │
│  │  │  │  └──────────────────────────────────────┘   │  │  │  │ │
│  │  │  └────────────────────────────────────────────────┘  │  │ │
│  │  └──────────────────────────────────────────────────────┘  │ │
│  │                            │                                │ │
│  │                            ▼                                │ │
│  │  ┌──────────────────────────────────────────────────────┐  │ │
│  │  │                  Backend API                          │  │ │
│  │  │  POST /api/v1/scorm/:id/initialize                    │  │ │
│  │  │  POST /api/v1/scorm/:id/get-value                     │  │ │
│  │  │  POST /api/v1/scorm/:id/set-value                     │  │ │
│  │  │  POST /api/v1/scorm/:id/commit                        │  │ │
│  │  │  POST /api/v1/scorm/:id/terminate                     │  │ │
│  │  └──────────────────────────────────────────────────────┘  │ │
│  └────────────────────────────────────────────────────────────┘ │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

#### Схема базы данных:

```sql
-- SCORM пакеты
CREATE TABLE scorm_packages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    
    -- Метаданные пакета
    title VARCHAR(500) NOT NULL,
    description TEXT,
    version VARCHAR(20) NOT NULL,        -- '1.2', '2004_3rd', '2004_4th'
    
    -- Файловая информация
    original_filename VARCHAR(255),
    storage_path VARCHAR(500) NOT NULL,   -- Путь в S3/MinIO
    manifest_path VARCHAR(500),           -- Путь к imsmanifest.xml
    launch_path VARCHAR(500) NOT NULL,    -- Путь к стартовому SCO
    
    -- Распарсенный манифест
    manifest_data JSONB,                  -- Полный манифест в JSON
    organizations JSONB,                  -- Структура курса
    
    -- Настройки
    mastery_score INTEGER,                -- Проходной балл
    max_time_allowed INTERVAL,            -- Лимит времени
    
    -- Статус
    status VARCHAR(20) DEFAULT 'active',  -- 'active', 'archived', 'draft'
    
    -- Метаданные
    uploaded_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- SCO (Sharable Content Objects) - отдельные модули внутри пакета
CREATE TABLE scorm_scos (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    package_id UUID NOT NULL REFERENCES scorm_packages(id) ON DELETE CASCADE,
    
    -- Идентификация
    identifier VARCHAR(255) NOT NULL,     -- ID из манифеста
    title VARCHAR(500) NOT NULL,
    
    -- Запуск
    launch_path VARCHAR(500) NOT NULL,    -- Относительный путь
    
    -- Позиция в структуре
    parent_id UUID REFERENCES scorm_scos(id),
    sort_order INTEGER DEFAULT 0,
    
    -- Настройки SCO
    mastery_score INTEGER,
    max_time_allowed INTERVAL,
    time_limit_action VARCHAR(50),        -- 'exit,no message', 'continue,message'
    
    -- Prerequisites (для SCORM 2004)
    prerequisites TEXT,                    -- Условия доступа
    
    -- Метаданные
    sco_type VARCHAR(20) DEFAULT 'sco',   -- 'sco', 'asset'
    
    created_at TIMESTAMP DEFAULT NOW()
);

-- Попытки прохождения SCORM
CREATE TABLE scorm_attempts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    
    -- Связи
    package_id UUID NOT NULL REFERENCES scorm_packages(id),
    sco_id UUID NOT NULL REFERENCES scorm_scos(id),
    user_id UUID NOT NULL REFERENCES users(id),
    
    -- Контекст (в каком курсе/модуле используется)
    course_id UUID REFERENCES courses(id),
    module_id UUID,
    
    -- Статус
    attempt_number INTEGER DEFAULT 1,
    status VARCHAR(50) DEFAULT 'not attempted',  -- SCORM cmi.core.lesson_status
    -- 'passed', 'completed', 'failed', 'incomplete', 'browsed', 'not attempted'
    
    -- Оценка
    score_raw DECIMAL(10,2),
    score_min DECIMAL(10,2),
    score_max DECIMAL(10,2),
    score_scaled DECIMAL(5,4),            -- Для SCORM 2004 (-1 to 1)
    
    -- Время
    session_time INTERVAL,                 -- Время текущей сессии
    total_time INTERVAL,                   -- Общее время
    
    -- Прогресс (SCORM 2004)
    progress_measure DECIMAL(5,4),         -- 0 to 1
    completion_status VARCHAR(50),         -- 'completed', 'incomplete', 'not attempted', 'unknown'
    success_status VARCHAR(50),            -- 'passed', 'failed', 'unknown'
    
    -- Данные для возобновления
    suspend_data TEXT,                     -- Сохраненное состояние
    location VARCHAR(1000),                -- cmi.core.lesson_location
    
    -- Entry/Exit
    entry VARCHAR(50),                     -- 'ab-initio', 'resume', ''
    exit_type VARCHAR(50),                 -- 'time-out', 'suspend', 'logout', 'normal', ''
    
    -- Временные метки
    started_at TIMESTAMP,
    last_accessed_at TIMESTAMP,
    completed_at TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(package_id, sco_id, user_id, attempt_number)
);

-- Runtime данные SCORM (все cmi.* значения)
CREATE TABLE scorm_runtime_data (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attempt_id UUID NOT NULL REFERENCES scorm_attempts(id) ON DELETE CASCADE,
    
    -- Ключ-значение
    element VARCHAR(255) NOT NULL,         -- cmi.core.student_name, cmi.interactions.0.id
    value TEXT,
    
    -- Для массивов (interactions, objectives)
    element_index INTEGER,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(attempt_id, element, element_index)
);

-- Взаимодействия (ответы на вопросы)
CREATE TABLE scorm_interactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attempt_id UUID NOT NULL REFERENCES scorm_attempts(id) ON DELETE CASCADE,
    
    -- Идентификация
    interaction_id VARCHAR(255) NOT NULL,  -- ID из контента
    interaction_index INTEGER NOT NULL,    -- Порядковый номер
    
    -- Тип взаимодействия
    interaction_type VARCHAR(50),          -- 'true-false', 'choice', 'fill-in', 'matching', 'performance', 'sequencing', 'likert', 'numeric', 'other'
    
    -- Вопрос и ответ
    description TEXT,                      -- Текст вопроса
    correct_responses TEXT[],              -- Правильные ответы
    learner_response TEXT,                 -- Ответ студента
    result VARCHAR(50),                    -- 'correct', 'incorrect', 'unanticipated', 'neutral'
    
    -- Оценка
    weighting DECIMAL(10,4),
    latency INTERVAL,                      -- Время на ответ
    
    -- Метки
    timestamp TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(attempt_id, interaction_index)
);

-- Objectives (цели обучения)
CREATE TABLE scorm_objectives (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    attempt_id UUID NOT NULL REFERENCES scorm_attempts(id) ON DELETE CASCADE,
    
    -- Идентификация
    objective_id VARCHAR(255) NOT NULL,
    objective_index INTEGER NOT NULL,
    
    -- Статус и оценка
    status VARCHAR(50),                    -- Аналогично lesson_status
    score_raw DECIMAL(10,2),
    score_min DECIMAL(10,2),
    score_max DECIMAL(10,2),
    score_scaled DECIMAL(5,4),
    
    -- Прогресс (SCORM 2004)
    progress_measure DECIMAL(5,4),
    completion_status VARCHAR(50),
    success_status VARCHAR(50),
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(attempt_id, objective_index)
);

-- Индексы
CREATE INDEX idx_scorm_packages_tenant ON scorm_packages(tenant_id);
CREATE INDEX idx_scorm_attempts_user ON scorm_attempts(user_id);
CREATE INDEX idx_scorm_attempts_package ON scorm_attempts(package_id);
CREATE INDEX idx_scorm_attempts_course ON scorm_attempts(course_id);
CREATE INDEX idx_scorm_runtime_attempt ON scorm_runtime_data(attempt_id);
CREATE INDEX idx_scorm_interactions_attempt ON scorm_interactions(attempt_id);
```

#### Frontend SCORM API Adapter:

```javascript
// frontend/src/lib/scorm/SCORMAdapter.js

/**
 * SCORM API Adapter
 * Реализует window.API (SCORM 1.2) и window.API_1484_11 (SCORM 2004)
 */
class SCORMAdapter {
  constructor(config) {
    this.attemptId = config.attemptId;
    this.apiEndpoint = config.apiEndpoint;
    this.version = config.version || '1.2';
    
    this.initialized = false;
    this.terminated = false;
    this.lastError = '0';
    
    // Кэш данных для уменьшения запросов
    this.dataCache = {};
    this.dirtyData = {};
    
    // Автосохранение
    this.autoCommitInterval = null;
    this.autoCommitDelay = 60000; // 1 минута
  }

  // ==================== SCORM 1.2 API ====================
  
  LMSInitialize(param) {
    if (this.initialized) {
      this.lastError = '101'; // Already initialized
      return 'false';
    }
    
    try {
      // Синхронный запрос для инициализации
      const response = this._apiCall('initialize', {});
      this.dataCache = response.data || {};
      this.initialized = true;
      this.lastError = '0';
      
      // Запускаем автосохранение
      this._startAutoCommit();
      
      return 'true';
    } catch (error) {
      this.lastError = '101';
      console.error('LMSInitialize failed:', error);
      return 'false';
    }
  }
  
  LMSFinish(param) {
    if (!this.initialized) {
      this.lastError = '301'; // Not initialized
      return 'false';
    }
    
    if (this.terminated) {
      this.lastError = '101'; // Already terminated
      return 'false';
    }
    
    try {
      // Сохраняем все несохраненные данные
      this.LMSCommit('');
      
      // Завершаем сессию
      this._apiCall('terminate', {});
      
      this.terminated = true;
      this._stopAutoCommit();
      this.lastError = '0';
      
      return 'true';
    } catch (error) {
      this.lastError = '101';
      console.error('LMSFinish failed:', error);
      return 'false';
    }
  }
  
  LMSGetValue(element) {
    if (!this.initialized || this.terminated) {
      this.lastError = '301';
      return '';
    }
    
    // Проверяем кэш
    if (this.dataCache.hasOwnProperty(element)) {
      this.lastError = '0';
      return this.dataCache[element];
    }
    
    // Read-only системные элементы
    const readOnlyElements = {
      'cmi.core._children': 'student_id,student_name,lesson_location,credit,lesson_status,entry,score,total_time,lesson_mode,exit,session_time',
      'cmi.core.score._children': 'raw,min,max',
      'cmi.student_data._children': 'mastery_score,max_time_allowed,time_limit_action',
    };
    
    if (readOnlyElements[element]) {
      this.lastError = '0';
      return readOnlyElements[element];
    }
    
    // Запрашиваем с сервера
    try {
      const response = this._apiCall('get-value', { element });
      this.dataCache[element] = response.value || '';
      this.lastError = '0';
      return this.dataCache[element];
    } catch (error) {
      this.lastError = '201'; // Invalid argument
      return '';
    }
  }
  
  LMSSetValue(element, value) {
    if (!this.initialized || this.terminated) {
      this.lastError = '301';
      return 'false';
    }
    
    // Проверяем read-only элементы
    const readOnlyElements = [
      'cmi.core.student_id',
      'cmi.core.student_name',
      'cmi.core.credit',
      'cmi.core.entry',
      'cmi.core.total_time',
      'cmi.core.lesson_mode',
      'cmi.core._children',
      'cmi.core.score._children',
    ];
    
    if (readOnlyElements.includes(element)) {
      this.lastError = '403'; // Read-only element
      return 'false';
    }
    
    // Валидация значений
    if (!this._validateValue(element, value)) {
      this.lastError = '405'; // Incorrect data type
      return 'false';
    }
    
    // Сохраняем в кэш и помечаем как "грязные"
    this.dataCache[element] = value;
    this.dirtyData[element] = value;
    this.lastError = '0';
    
    return 'true';
  }
  
  LMSCommit(param) {
    if (!this.initialized || this.terminated) {
      this.lastError = '301';
      return 'false';
    }
    
    if (Object.keys(this.dirtyData).length === 0) {
      this.lastError = '0';
      return 'true';
    }
    
    try {
      this._apiCall('commit', { data: this.dirtyData });
      this.dirtyData = {};
      this.lastError = '0';
      return 'true';
    } catch (error) {
      this.lastError = '101';
      console.error('LMSCommit failed:', error);
      return 'false';
    }
  }
  
  LMSGetLastError() {
    return this.lastError;
  }
  
  LMSGetErrorString(errorCode) {
    const errors = {
      '0': 'No Error',
      '101': 'General Exception',
      '201': 'Invalid Argument Error',
      '202': 'Element Cannot Have Children',
      '203': 'Element Not An Array',
      '301': 'Not Initialized',
      '401': 'Not Implemented Error',
      '402': 'Invalid Set Value',
      '403': 'Element Is Read Only',
      '404': 'Element Is Write Only',
      '405': 'Incorrect Data Type',
    };
    return errors[errorCode] || 'Unknown Error';
  }
  
  LMSGetDiagnostic(errorCode) {
    return this.LMSGetErrorString(errorCode);
  }
  
  // ==================== SCORM 2004 API ====================
  // (Аналогичные методы с префиксом без LMS)
  
  Initialize(param) { return this.LMSInitialize(param); }
  Terminate(param) { return this.LMSFinish(param); }
  GetValue(element) { return this.LMSGetValue(this._convert2004To12(element)); }
  SetValue(element, value) { return this.LMSSetValue(this._convert2004To12(element), value); }
  Commit(param) { return this.LMSCommit(param); }
  GetLastError() { return this.LMSGetLastError(); }
  GetErrorString(errorCode) { return this.LMSGetErrorString(errorCode); }
  GetDiagnostic(errorCode) { return this.LMSGetDiagnostic(errorCode); }
  
  // ==================== Private Methods ====================
  
  _apiCall(method, data) {
    const xhr = new XMLHttpRequest();
    xhr.open('POST', `${this.apiEndpoint}/${method}`, false); // Синхронный для SCORM
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.send(JSON.stringify({
      attemptId: this.attemptId,
      ...data
    }));
    
    if (xhr.status !== 200) {
      throw new Error(`API call failed: ${xhr.status}`);
    }
    
    return JSON.parse(xhr.responseText);
  }
  
  _validateValue(element, value) {
    // Валидация по типу элемента
    const validators = {
      'cmi.core.lesson_status': (v) => ['passed', 'completed', 'failed', 'incomplete', 'browsed', 'not attempted'].includes(v),
      'cmi.core.score.raw': (v) => !isNaN(parseFloat(v)),
      'cmi.core.score.min': (v) => !isNaN(parseFloat(v)),
      'cmi.core.score.max': (v) => !isNaN(parseFloat(v)),
      'cmi.core.session_time': (v) => /^\d{2,4}:\d{2}:\d{2}(\.\d{1,2})?$/.test(v),
    };
    
    if (validators[element]) {
      return validators[element](value);
    }
    
    return true;
  }
  
  _convert2004To12(element) {
    // Маппинг SCORM 2004 -> SCORM 1.2
    const mapping = {
      'cmi.completion_status': 'cmi.core.lesson_status',
      'cmi.success_status': 'cmi.core.lesson_status',
      'cmi.score.scaled': 'cmi.core.score.raw',
      'cmi.score.raw': 'cmi.core.score.raw',
      'cmi.location': 'cmi.core.lesson_location',
      'cmi.exit': 'cmi.core.exit',
    };
    return mapping[element] || element;
  }
  
  _startAutoCommit() {
    this.autoCommitInterval = setInterval(() => {
      if (Object.keys(this.dirtyData).length > 0) {
        this.LMSCommit('');
      }
    }, this.autoCommitDelay);
  }
  
  _stopAutoCommit() {
    if (this.autoCommitInterval) {
      clearInterval(this.autoCommitInterval);
      this.autoCommitInterval = null;
    }
  }
}

// Регистрируем глобально для доступа из iframe
window.SCORMAdapter = SCORMAdapter;
```

```typescript
// frontend/src/components/SCORMPlayer.tsx
import React, { useEffect, useRef, useState } from 'react';

interface SCORMPlayerProps {
  packageId: string;
  attemptId: string;
  launchUrl: string;
  title: string;
  onComplete?: (data: { status: string; score?: number }) => void;
}

export function SCORMPlayer({ packageId, attemptId, launchUrl, title, onComplete }: SCORMPlayerProps) {
  const iframeRef = useRef<HTMLIFrameElement>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    // Инициализируем SCORM API до загрузки контента
    const adapter = new window.SCORMAdapter({
      attemptId,
      apiEndpoint: `/api/v1/scorm/${packageId}`,
      version: '1.2',
    });

    // Регистрируем API глобально (SCORM ищет window.API)
    window.API = adapter;
    window.API_1484_11 = adapter; // Для SCORM 2004

    // Слушаем события от iframe
    const handleMessage = (event: MessageEvent) => {
      if (event.data.type === 'scorm_complete') {
        onComplete?.(event.data);
      }
    };
    window.addEventListener('message', handleMessage);

    // Очистка при размонтировании
    return () => {
      window.removeEventListener('message', handleMessage);
      // Завершаем сессию если не завершена
      if (window.API && !window.API.terminated) {
        window.API.LMSFinish('');
      }
      delete window.API;
      delete window.API_1484_11;
    };
  }, [attemptId, packageId, onComplete]);

  const handleIframeLoad = () => {
    setLoading(false);
  };

  const handleIframeError = () => {
    setError('Failed to load SCORM content');
    setLoading(false);
  };

  return (
    <div className="scorm-player h-full flex flex-col">
      <div className="scorm-header bg-gray-100 p-2 flex items-center justify-between">
        <h3 className="font-medium">{title}</h3>
        <button 
          onClick={() => window.API?.LMSCommit('')}
          className="text-sm text-blue-600 hover:text-blue-800"
        >
          Save Progress
        </button>
      </div>
      
      {loading && (
        <div className="flex-1 flex items-center justify-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
        </div>
      )}
      
      {error && (
        <div className="flex-1 flex items-center justify-center text-red-600">
          {error}
        </div>
      )}
      
      <iframe
        ref={iframeRef}
        src={launchUrl}
        className={`flex-1 w-full border-0 ${loading ? 'hidden' : ''}`}
        onLoad={handleIframeLoad}
        onError={handleIframeError}
        sandbox="allow-scripts allow-same-origin allow-forms allow-popups"
        title={title}
      />
    </div>
  );
}
```

---

### 4.6 Сложность освоения и реализации

#### Оценка сложности:

| Аспект | Сложность | Комментарий |
|--------|-----------|-------------|
| **Парсинг манифеста** | 🟡 Средняя | XML парсинг, множество вариаций |
| **SCORM 1.2 RTE** | 🟡 Средняя | Ограниченный API, хорошо документирован |
| **SCORM 2004 RTE** | 🔴 Высокая | Sequencing, более сложная модель |
| **Frontend Player** | 🟡 Средняя | iframe, кроссбраузерность |
| **Хранение данных** | 🟢 Низкая | Простая модель данных |
| **Тестирование** | 🔴 Высокая | Много разных пакетов, edge cases |

#### Временные оценки:

```
SCORM 1.2 (базовая поддержка):
├── Парсинг манифеста: 3-4 дня
├── Upload & storage: 2-3 дня
├── RTE API backend: 3-4 дня
├── Frontend player: 4-5 дня
├── Тестирование: 3-4 дня
└── Итого: 3 недели

SCORM 2004 (расширенная поддержка):
├── Расширение парсинга: 2-3 дня
├── Sequencing engine: 5-7 дней
├── Расширение RTE API: 3-4 дня
├── Тестирование: 4-5 дней
└── Итого: +2 недели

Общее время: 5-6 недель (один разработчик)
```

#### Типичные проблемы и решения:

| Проблема | Причина | Решение |
|----------|---------|---------|
| API не найден | iframe загружается раньше API | Инициализировать API до iframe |
| Синхронные вызовы | SCORM требует синхронный API | XMLHttpRequest синхронный режим |
| Cross-origin | SCORM контент на другом домене | Одинаковый origin или postMessage |
| Кодировка | UTF-8 vs Windows-1251 | Детектировать и конвертировать |
| Некорректный манифест | Разные authoring tools | Гибкий парсер с fallbacks |

#### Тестовые SCORM-пакеты:

```
Рекомендуемые пакеты для тестирования:

ADL (официальные тесты):
• SCORM 1.2 Test Suite
• SCORM 2004 Conformance Test Suite

Open Source:
• Rustici Golf Examples (разные сценарии)
• SCORM Cloud Test Packages

Создать самим:
• Articulate Storyline (trial)
• Adobe Captivate (trial)
• iSpring Suite (trial)
```

---

### 4.7 Источники для дальнейшего изучения

#### Официальные спецификации:

| Ресурс | Ссылка | Описание |
|--------|--------|----------|
| **ADL SCORM** | [adlnet.gov/scorm](https://adlnet.gov/projects/scorm/) | Официальная документация |
| **SCORM 1.2 RTE** | [scorm.com/scorm-explained/scorm-12](https://scorm.com/scorm-explained/technical-scorm/run-time/) | Run-Time Environment |
| **SCORM 2004 Spec** | [adlnet.gov/scorm/scorm-2004-4th](https://adlnet.gov/projects/scorm/scorm-2004-4th-edition/) | 4th Edition |

#### Инструменты:

| Инструмент | Описание | Ссылка |
|------------|----------|--------|
| **SCORM Cloud** | Тестирование пакетов | [cloud.scorm.com](https://cloud.scorm.com/) |
| **Rustici Engine** | Production SCORM engine | [rusticisoftware.com](https://rusticisoftware.com/) |
| **ADL Test Suite** | Conformance testing | [adlnet.gov](https://adlnet.gov/) |

#### Библиотеки (JavaScript):

| Библиотека | Описание | Ссылка |
|------------|----------|--------|
| **pipwerks SCORM** | JavaScript SCORM wrapper | [pipwerks.com](https://pipwerks.com/laboratory/scorm/) |
| **scorm-again** | Modern SCORM wrapper | [npm: scorm-again](https://www.npmjs.com/package/scorm-again) |

#### Authoring Tools:

| Инструмент | Тип | Популярность |
|------------|-----|--------------|
| **Articulate Storyline** | Desktop | #1 в индустрии |
| **Adobe Captivate** | Desktop | Сильный в симуляциях |
| **iSpring Suite** | PowerPoint add-in | Простой в использовании |
| **Lectora** | Desktop | Enterprise features |
| **Adapt Framework** | Open source, web | Бесплатный |

---

### 4.8 Чек-лист реализации

```
Phase 1: Package Management (Day 1-5)
□ API для upload ZIP
□ Распаковка и валидация
□ Парсинг imsmanifest.xml
□ Сохранение в S3/MinIO
□ Database schema
□ Список пакетов API
□ Удаление пакетов

Phase 2: SCORM 1.2 Runtime (Day 6-12)
□ LMSInitialize endpoint
□ LMSGetValue endpoint
□ LMSSetValue endpoint
□ LMSCommit endpoint
□ LMSFinish endpoint
□ Error handling
□ Data validation
□ Session management

Phase 3: Frontend Player (Day 13-18)
□ SCORMAdapter JavaScript
□ Player component (React)
□ Navigation (TOC)
□ Progress tracking UI
□ Error handling UI
□ Fullscreen mode

Phase 4: Integration (Day 19-22)
□ Embedding в Course module
□ Grade passback
□ Completion tracking
□ Reports API

Phase 5: SCORM 2004 (Day 23-30)
□ Extended data model
□ Sequencing (basic)
□ Navigation controls
□ Testing с 2004 пакетами

Phase 6: Testing & Polish (Day 31-35)
□ Тестирование с Articulate
□ Тестирование с Captivate
□ Тестирование с iSpring
□ ADL Test Suite
□ Edge cases
□ Документация
```

---

## 5. Video Conferencing Integration (Видеоконференции)

### 5.1 Определение

**Video Conferencing Integration** — это встраивание видеозвонков и онлайн-встреч непосредственно в образовательную платформу. Позволяет проводить живые лекции, семинары, консультации и защиты без переключения между приложениями.

**Ключевая идея:** Студенты и преподаватели должны иметь возможность начать видеозвонок в один клик прямо из интерфейса LMS, без необходимости отдельно логиниться в Zoom/Teams/Meet.

#### Типы интеграции:

| Тип | Описание | Примеры | Сложность |
|-----|----------|---------|-----------|
| **Redirect** | Перенаправление на внешний сервис | Ссылка на Zoom | 🟢 Простая |
| **Embed** | Встраивание через iframe | Jitsi Meet | 🟡 Средняя |
| **API Integration** | Программное управление через API | Zoom API, Teams Graph API | 🔴 Высокая |
| **Self-hosted** | Собственный сервер видео | BigBlueButton, Jitsi | 🔴 Высокая |

#### Популярные решения:

| Платформа | Тип | Стоимость | Особенности |
|-----------|-----|-----------|-------------|
| **Zoom** | SaaS | $14.99+/host/месяц | #1 по популярности, отличный API |
| **Microsoft Teams** | SaaS | Включен в M365 | Интеграция с Office, Graph API |
| **Google Meet** | SaaS | Включен в Workspace | Простота, Calendar интеграция |
| **BigBlueButton** | Open Source | Бесплатно | Создан для обучения |
| **Jitsi Meet** | Open Source | Бесплатно | Простой, без регистрации |
| **Whereby** | SaaS | $6.99+/месяц | Embed-first, красивый UI |
| **Daily.co** | SaaS/API | Pay-per-minute | Отличный API для разработчиков |

#### Архитектура видеоконференций:

```
┌─────────────────────────────────────────────────────────────────┐
│                Video Conferencing Architecture                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                    Our LMS Platform                      │    │
│  │                                                          │    │
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────┐  │    │
│  │  │   Course    │  │  Schedule   │  │    Meeting      │  │    │
│  │  │   Module    │  │  Calendar   │  │    History      │  │    │
│  │  └──────┬──────┘  └──────┬──────┘  └────────┬────────┘  │    │
│  │         │                │                  │           │    │
│  │         └────────────────┼──────────────────┘           │    │
│  │                          ▼                              │    │
│  │  ┌─────────────────────────────────────────────────────┐│    │
│  │  │              Video Service Abstraction              ││    │
│  │  │  ┌──────────────────────────────────────────────┐  ││    │
│  │  │  │ interface VideoProvider {                     │  ││    │
│  │  │  │   CreateMeeting(opts) -> Meeting             │  ││    │
│  │  │  │   GetMeetingInfo(id) -> Meeting              │  ││    │
│  │  │  │   UpdateMeeting(id, opts) -> Meeting         │  ││    │
│  │  │  │   DeleteMeeting(id) -> bool                  │  ││    │
│  │  │  │   GetJoinURL(id, user) -> string             │  ││    │
│  │  │  │   GetRecordings(id) -> []Recording           │  ││    │
│  │  │  │ }                                            │  ││    │
│  │  │  └──────────────────────────────────────────────┘  ││    │
│  │  └─────────────────────────────────────────────────────┘│    │
│  │                          │                              │    │
│  └──────────────────────────┼──────────────────────────────┘    │
│                             │                                    │
│           ┌─────────────────┼─────────────────┐                 │
│           │                 │                 │                 │
│           ▼                 ▼                 ▼                 │
│  ┌─────────────┐   ┌─────────────┐   ┌─────────────────┐       │
│  │    Zoom     │   │   Teams     │   │  BigBlueButton  │       │
│  │    API      │   │  Graph API  │   │     API         │       │
│  └─────────────┘   └─────────────┘   └─────────────────┘       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

#### Сценарии использования в образовании:

| Сценарий | Описание | Участники | Частота |
|----------|----------|-----------|---------|
| **Лекция** | Преподаватель ведет, студенты слушают | 1 → 100+ | Регулярно |
| **Семинар** | Обсуждение, работа в группах | 10-30 | Регулярно |
| **Консультация** | 1:1 или малая группа | 2-5 | По запросу |
| **Защита** | Презентация + комиссия | 5-10 | Периодически |
| **Office Hours** | Приемные часы преподавателя | 1-5 | Регулярно |
| **Вебинар** | Массовое мероприятие | 100-1000+ | Редко |

---

### 5.2 Почему важно реализовать

#### Бизнес-причины:

| Причина | Описание | Влияние |
|---------|----------|---------|
| **Post-COVID стандарт** | Онлайн/гибрид стал нормой | Must-have для 2024+ |
| **Конкурентное преимущество** | Бесшовный UX | Удержание пользователей |
| **Университетские RFP** | Требование в тендерах | Доступ к крупным контрактам |
| **Экономия времени** | Не нужно переключаться | +15% продуктивности |
| **Централизация данных** | Записи, attendance в LMS | Единый источник правды |

#### Образовательные причины:

```
📊 Статистика онлайн-обучения (2024):

• 70% университетов используют гибридный формат
• 85% студентов ожидают опцию онлайн-участия
• 60% курсов имеют хотя бы одну онлайн-сессию
• Записи лекций повышают успеваемость на 12%
• Онлайн office hours увеличивают посещаемость на 40%

💡 Новые реалии:
• Студенты с работой → нужна гибкость
• Международные студенты → разные часовые пояса
• Студенты с ограниченными возможностями → доступность
• PhD защиты → часто с удаленными экспертами
```

#### Технические причины:

1. **Single Sign-On** — один вход для LMS и видео
2. **Автоматическое создание встреч** — при создании занятия
3. **Attendance tracking** — автоматическое отслеживание присутствия
4. **Записи в контексте курса** — доступ к записям в нужном месте
5. **Интеграция с журналом** — участие влияет на оценку

#### Сравнение: с интеграцией vs без интеграции

```
Без интеграции (типичный workflow):
├─ Преподаватель создает встречу в Zoom отдельно
├─ Копирует ссылку
├─ Вставляет в объявление в LMS
├─ Студенты переходят по ссылке
├─ Логинятся в Zoom отдельно
├─ Записи остаются в Zoom (не в LMS)
├─ Attendance считают вручную
└─ Время: 10+ минут настройки, данные разрозненны

С интеграцией:
├─ Преподаватель нажимает "Создать онлайн-занятие"
├─ Встреча создается автоматически (Zoom/Teams/BBB)
├─ Студенты нажимают "Присоединиться" в LMS
├─ SSO — вход автоматический
├─ Записи сохраняются в курс
├─ Attendance автоматически в журнале
└─ Время: 1 минута, все данные централизованы ✅
```

---

### 5.3 Что дает конечному пользователю

#### Для студентов:

| Функция | Описание | Польза |
|---------|----------|--------|
| **One-click join** | Присоединиться в один клик | Экономия времени |
| **Расписание в одном месте** | Все занятия в календаре LMS | Не пропустить встречу |
| **Записи в курсе** | Записи рядом с материалами | Удобный пересмотр |
| **Мобильный доступ** | Присоединиться с телефона | Гибкость |
| **Чат интегрирован** | Вопросы сохраняются | Контекст обсуждений |

#### Для преподавателей:

| Функция | Описание | Польза |
|---------|----------|--------|
| **Автоматическое создание** | Встреча создается с курсом | Меньше рутины |
| **Attendance автоматически** | Кто был — видно в журнале | Экономия 30 мин/занятие |
| **Breakout rooms** | Работа в малых группах | Интерактивное обучение |
| **Whiteboard** | Совместная доска | Визуальное объяснение |
| **Recording management** | Управление записями | Контроль контента |
| **Polling/Quiz** | Опросы во время занятия | Вовлечение |

#### Для администраторов:

| Функция | Описание | Польза |
|---------|----------|--------|
| **Единый биллинг** | Расходы на видео в LMS | Прозрачность |
| **Usage analytics** | Статистика использования | ROI видео |
| **Compliance** | Записи для аудита | Регуляторные требования |
| **Central management** | Настройки для всей организации | Контроль |

#### Конкретные сценарии:

```
Сценарий 1: Регулярная лекция
├─ Преподаватель: создает курс с расписанием
├─ Система: автоматически создает Zoom для каждого занятия
├─ За 15 мин: напоминание студентам
├─ Студенты: нажимают "Join" в LMS
├─ После: запись автоматически доступна в курсе
├─ Attendance: автоматически в журнале
└─ Студенты, пропустившие: смотрят запись ✓

Сценарий 2: PhD защита
├─ Секретарь: создает "Защита диссертации Иванова"
├─ Приглашает: научрука, оппонентов, комиссию
├─ Внешний эксперт: получает ссылку без регистрации
├─ Защита: проходит онлайн/гибрид
├─ Запись: архивируется для документации
└─ Протокол: генерируется с attendance ✓

Сценарий 3: Office Hours (консультации)
├─ Преподаватель: устанавливает слоты в календаре
├─ Студенты: бронируют удобное время
├─ Система: создает приватную встречу
├─ 1:1 консультация: проходит в видео
├─ Заметки: преподаватель пишет в карточку студента
└─ История: все консультации логируются ✓

Сценарий 4: Гостевая лекция
├─ Преподаватель: приглашает эксперта из индустрии
├─ Эксперт: получает временный доступ
├─ Студенты: присоединяются как обычно
├─ Q&A: через чат, модерирует преподаватель
├─ Запись: доступна студентам
└─ Эксперт: доступ автоматически закрывается ✓
```

---

### 5.4 Какие пользователи предпочитают эту функцию

#### Целевые сегменты:

| Сегмент | Важность видео | Предпочтения |
|---------|---------------|--------------|
| **Университеты** | Критическая | Zoom, Teams, BBB |
| **Корпоративное обучение** | Высокая | Teams, Zoom |
| **K-12** | Высокая | Google Meet, Zoom |
| **Курсы/Bootcamps** | Критическая | Zoom, собственное решение |
| **Репетиторство** | Критическая | 1:1 видео |

#### Типичные требования по сегментам:

```
🎓 Университеты:
• Интеграция с существующим Zoom/Teams
• Поддержка 500+ участников (лекции)
• Записи с транскрипцией
• Attendance для администрации
• Breakout rooms для семинаров

💼 Корпорации:
• Microsoft Teams (уже оплачен)
• Соответствие безопасности (SOC2, GDPR)
• Интеграция с HR системами
• Отчетность по обучению

🏫 K-12:
• Простота для учителей
• Родительский контроль
• Безопасность детей
• Google Workspace интеграция

🚀 Bootcamps:
• Масштабируемость (пиковые нагрузки)
• Запись и воспроизведение
• Breakout для pair programming
• Screen sharing качество
```

#### Вопросы при продаже:

```
Типичные вопросы от клиентов:

1. "Поддерживаете ли вы Zoom/Teams?"
2. "Можно ли использовать наш существующий Zoom аккаунт?"
3. "Как работает SSO для видео?"
4. "Где хранятся записи?"
5. "Автоматический ли attendance?"
6. "Есть ли лимиты на участников/время?"
7. "Поддерживается ли BigBlueButton?" (open source)
8. "Можно ли проводить вебинары?"

Без видеоинтеграции = "Нет" = потеря 80%+ образовательных клиентов
```

---

### 5.5 Как интегрировать в наше приложение

#### Выбор стратегии интеграции:

| Подход | Плюсы | Минусы | Рекомендация |
|--------|-------|--------|--------------|
| **Zoom API** | Популярность, отличный API | Платный, зависимость | ✅ Основной |
| **BigBlueButton** | Бесплатный, для образования | Нужен сервер, UI старый | ✅ Open source опция |
| **Jitsi** | Простой, бесплатный | Ограничения масштаба | 🔶 Для малых групп |
| **Daily.co** | Отличный API, embed | Дорого при масштабе | 🔶 Для premium |
| **Teams Graph API** | Интеграция M365 | Сложный API | 🔶 Для Enterprise |

#### Рекомендуемая архитектура:

```
┌─────────────────────────────────────────────────────────────────┐
│                 Recommended Architecture                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Primary: Zoom (for paid customers with Zoom licenses)          │
│  Secondary: BigBlueButton (for self-hosted / budget-conscious)  │
│  Fallback: Jitsi (for quick 1:1 without setup)                  │
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                Video Provider Interface                  │    │
│  │                                                          │    │
│  │  CreateMeeting(title, start, duration, opts) → Meeting  │    │
│  │  GetMeeting(id) → Meeting                               │    │
│  │  UpdateMeeting(id, opts) → Meeting                      │    │
│  │  DeleteMeeting(id) → bool                               │    │
│  │  GetJoinURL(meetingId, userId, role) → JoinInfo         │    │
│  │  EndMeeting(id) → bool                                  │    │
│  │  GetRecordings(meetingId) → []Recording                 │    │
│  │  GetParticipants(meetingId) → []Participant             │    │
│  │  RegisterWebhook(url, events) → bool                    │    │
│  │                                                          │    │
│  └─────────────────────────────────────────────────────────┘    │
│           │                │                │                    │
│           ▼                ▼                ▼                    │
│  ┌─────────────┐  ┌─────────────────┐  ┌─────────────┐          │
│  │ ZoomProvider│  │ BBBProvider     │  │JitsiProvider│          │
│  │             │  │                 │  │             │          │
│  │ • OAuth 2.0 │  │ • API Secret    │  │ • JWT Token │          │
│  │ • REST API  │  │ • Checksum auth │  │ • REST API  │          │
│  │ • Webhooks  │  │ • Callbacks     │  │ • Simple    │          │
│  └─────────────┘  └─────────────────┘  └─────────────┘          │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

#### Схема базы данных:

```sql
-- Настройки видеопровайдера для организации
CREATE TABLE video_provider_configs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    
    -- Провайдер
    provider VARCHAR(50) NOT NULL,        -- 'zoom', 'bbb', 'teams', 'jitsi'
    is_default BOOLEAN DEFAULT false,
    is_enabled BOOLEAN DEFAULT true,
    
    -- Zoom OAuth credentials
    zoom_account_id VARCHAR(255),
    zoom_client_id VARCHAR(255),
    zoom_client_secret VARCHAR(500),      -- Encrypted
    zoom_webhook_secret VARCHAR(255),
    
    -- BigBlueButton credentials
    bbb_server_url VARCHAR(500),
    bbb_shared_secret VARCHAR(255),       -- Encrypted
    
    -- Teams/Graph API
    teams_tenant_id VARCHAR(255),
    teams_client_id VARCHAR(255),
    teams_client_secret VARCHAR(500),     -- Encrypted
    
    -- Jitsi (self-hosted)
    jitsi_server_url VARCHAR(500),
    jitsi_app_id VARCHAR(255),
    jitsi_app_secret VARCHAR(255),        -- Encrypted for JWT
    
    -- Настройки по умолчанию
    default_settings JSONB DEFAULT '{
        "waiting_room": true,
        "mute_on_entry": true,
        "auto_recording": false,
        "recording_location": "cloud"
    }',
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(tenant_id, provider)
);

-- Видеовстречи
CREATE TABLE video_meetings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    
    -- Связь с контекстом
    course_id UUID REFERENCES courses(id),
    module_id UUID,
    schedule_item_id UUID,                -- Связь с расписанием
    
    -- Провайдер и внешний ID
    provider VARCHAR(50) NOT NULL,
    external_id VARCHAR(255),             -- ID встречи у провайдера
    
    -- Метаданные встречи
    title VARCHAR(500) NOT NULL,
    description TEXT,
    meeting_type VARCHAR(50) DEFAULT 'scheduled',  -- 'instant', 'scheduled', 'recurring'
    
    -- Время
    scheduled_start TIMESTAMP,
    scheduled_end TIMESTAMP,
    duration_minutes INTEGER,
    timezone VARCHAR(50) DEFAULT 'UTC',
    
    -- Recurrence (для регулярных занятий)
    recurrence_rule TEXT,                 -- RRULE format
    recurrence_parent_id UUID REFERENCES video_meetings(id),
    
    -- URLs
    join_url VARCHAR(1000),
    host_url VARCHAR(1000),               -- URL с правами хоста
    password VARCHAR(50),
    
    -- Настройки
    settings JSONB DEFAULT '{}',
    /*
    {
        "waiting_room": true,
        "mute_on_entry": true,
        "allow_screen_share": "host_only",
        "recording": "none", // "cloud", "local", "none"
        "breakout_rooms": false,
        "max_participants": 100
    }
    */
    
    -- Статус
    status VARCHAR(50) DEFAULT 'scheduled',  -- 'scheduled', 'started', 'ended', 'cancelled'
    started_at TIMESTAMP,
    ended_at TIMESTAMP,
    actual_duration_minutes INTEGER,
    
    -- Создатель
    created_by UUID REFERENCES users(id),
    host_user_id UUID REFERENCES users(id),
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Участники встречи
CREATE TABLE video_meeting_participants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    meeting_id UUID NOT NULL REFERENCES video_meetings(id) ON DELETE CASCADE,
    
    -- Участник
    user_id UUID REFERENCES users(id),    -- NULL для гостей
    external_email VARCHAR(255),          -- Для внешних участников
    display_name VARCHAR(255),
    
    -- Роль
    role VARCHAR(50) DEFAULT 'attendee',  -- 'host', 'co-host', 'attendee', 'panelist'
    
    -- Приглашение
    invitation_status VARCHAR(50) DEFAULT 'pending',  -- 'pending', 'accepted', 'declined'
    invited_at TIMESTAMP DEFAULT NOW(),
    
    -- Фактическое участие (заполняется webhook-ами)
    join_time TIMESTAMP,
    leave_time TIMESTAMP,
    duration_seconds INTEGER,
    
    -- Дополнительная информация
    device_type VARCHAR(50),              -- 'desktop', 'mobile', 'phone'
    connection_type VARCHAR(50),          -- 'wifi', 'cellular'
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(meeting_id, user_id)
);

-- Записи встреч
CREATE TABLE video_recordings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    meeting_id UUID NOT NULL REFERENCES video_meetings(id) ON DELETE CASCADE,
    
    -- Внешние данные
    external_id VARCHAR(255),             -- ID записи у провайдера
    provider VARCHAR(50) NOT NULL,
    
    -- Файл
    recording_type VARCHAR(50),           -- 'video', 'audio', 'transcript', 'chat'
    file_type VARCHAR(20),                -- 'mp4', 'mp3', 'vtt', 'txt'
    
    -- Локация
    storage_location VARCHAR(50),         -- 'provider_cloud', 's3', 'local'
    provider_url VARCHAR(1000),           -- URL у провайдера
    local_path VARCHAR(500),              -- Путь в S3/локальном хранилище
    
    -- Метаданные
    file_size_bytes BIGINT,
    duration_seconds INTEGER,
    
    -- Статус
    status VARCHAR(50) DEFAULT 'processing',  -- 'processing', 'ready', 'failed', 'deleted'
    processing_progress INTEGER,          -- 0-100
    
    -- Доступ
    is_public BOOLEAN DEFAULT false,
    password VARCHAR(50),
    expires_at TIMESTAMP,
    
    -- Транскрипция
    transcript TEXT,
    transcript_language VARCHAR(10),
    
    -- Метки времени
    recording_start TIMESTAMP,
    recording_end TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Attendance log (детальный лог присутствия)
CREATE TABLE video_attendance_log (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    meeting_id UUID NOT NULL REFERENCES video_meetings(id) ON DELETE CASCADE,
    participant_id UUID REFERENCES video_meeting_participants(id),
    
    -- Событие
    event_type VARCHAR(50) NOT NULL,      -- 'join', 'leave', 'mute', 'unmute', 'screen_share_start', etc.
    event_time TIMESTAMP NOT NULL,
    
    -- Детали
    details JSONB,
    
    created_at TIMESTAMP DEFAULT NOW()
);

-- Индексы
CREATE INDEX idx_video_meetings_tenant ON video_meetings(tenant_id);
CREATE INDEX idx_video_meetings_course ON video_meetings(course_id);
CREATE INDEX idx_video_meetings_start ON video_meetings(scheduled_start);
CREATE INDEX idx_video_meetings_status ON video_meetings(status);
CREATE INDEX idx_video_participants_meeting ON video_meeting_participants(meeting_id);
CREATE INDEX idx_video_participants_user ON video_meeting_participants(user_id);
CREATE INDEX idx_video_recordings_meeting ON video_recordings(meeting_id);
CREATE INDEX idx_video_attendance_meeting ON video_attendance_log(meeting_id);
```

#### Frontend компоненты:

```typescript
// frontend/src/components/VideoMeeting/MeetingButton.tsx
import React from 'react';
import { Video, ExternalLink } from 'lucide-react';
import { Button } from '@/components/ui/button';

interface MeetingButtonProps {
  meeting: {
    id: string;
    title: string;
    status: 'scheduled' | 'started' | 'ended';
    joinUrl: string;
    scheduledStart: string;
  };
  userRole: 'host' | 'attendee';
}

export function MeetingButton({ meeting, userRole }: MeetingButtonProps) {
  const isLive = meeting.status === 'started';
  const canStart = userRole === 'host' && meeting.status === 'scheduled';
  const canJoin = meeting.status === 'started' || 
                  (meeting.status === 'scheduled' && isWithin15MinutesOfStart(meeting.scheduledStart));

  const handleJoin = async () => {
    // Получаем персонализированную ссылку с SSO
    const response = await fetch(`/api/v1/meetings/${meeting.id}/join`, {
      method: 'POST',
    });
    const { joinUrl } = await response.json();
    
    // Открываем в новом окне
    window.open(joinUrl, '_blank', 'noopener,noreferrer');
  };

  const handleStart = async () => {
    await fetch(`/api/v1/meetings/${meeting.id}/start`, { method: 'POST' });
    handleJoin();
  };

  if (meeting.status === 'ended') {
    return (
      <Button variant="outline" disabled>
        <Video className="w-4 h-4 mr-2" />
        Meeting Ended
      </Button>
    );
  }

  if (canStart) {
    return (
      <Button onClick={handleStart} className="bg-green-600 hover:bg-green-700">
        <Video className="w-4 h-4 mr-2" />
        Start Meeting
      </Button>
    );
  }

  if (canJoin) {
    return (
      <Button onClick={handleJoin} className={isLive ? 'bg-red-600 hover:bg-red-700 animate-pulse' : ''}>
        <Video className="w-4 h-4 mr-2" />
        {isLive ? 'Join Live' : 'Join Meeting'}
        <ExternalLink className="w-3 h-3 ml-1" />
      </Button>
    );
  }

  return (
    <Button variant="outline" disabled>
      <Video className="w-4 h-4 mr-2" />
      Starts at {formatTime(meeting.scheduledStart)}
    </Button>
  );
}

function isWithin15MinutesOfStart(startTime: string): boolean {
  const start = new Date(startTime);
  const now = new Date();
  const diffMinutes = (start.getTime() - now.getTime()) / (1000 * 60);
  return diffMinutes <= 15 && diffMinutes >= -60; // 15 мин до, 60 мин после
}

function formatTime(dateString: string): string {
  return new Date(dateString).toLocaleTimeString('ru-RU', {
    hour: '2-digit',
    minute: '2-digit',
  });
}
```

```typescript
// frontend/src/components/VideoMeeting/CreateMeetingDialog.tsx
import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';

interface CreateMeetingForm {
  title: string;
  description?: string;
  scheduledStart: string;
  durationMinutes: number;
  waitingRoom: boolean;
  muteOnEntry: boolean;
  autoRecording: boolean;
}

interface CreateMeetingDialogProps {
  open: boolean;
  onClose: () => void;
  courseId: string;
  onCreated: (meeting: any) => void;
}

export function CreateMeetingDialog({ open, onClose, courseId, onCreated }: CreateMeetingDialogProps) {
  const [loading, setLoading] = useState(false);
  const { register, handleSubmit, watch } = useForm<CreateMeetingForm>({
    defaultValues: {
      durationMinutes: 60,
      waitingRoom: true,
      muteOnEntry: true,
      autoRecording: false,
    },
  });

  const onSubmit = async (data: CreateMeetingForm) => {
    setLoading(true);
    try {
      const response = await fetch('/api/v1/meetings', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          courseId,
          title: data.title,
          description: data.description,
          scheduledStart: new Date(data.scheduledStart).toISOString(),
          durationMinutes: data.durationMinutes,
          settings: {
            waitingRoom: data.waitingRoom,
            muteOnEntry: data.muteOnEntry,
            recording: data.autoRecording ? 'cloud' : 'none',
          },
        }),
      });
      
      const meeting = await response.json();
      onCreated(meeting);
      onClose();
    } catch (error) {
      console.error('Failed to create meeting:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[500px]">
        <DialogHeader>
          <DialogTitle>Schedule Video Meeting</DialogTitle>
        </DialogHeader>
        
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <Label htmlFor="title">Meeting Title</Label>
            <Input
              id="title"
              {...register('title', { required: true })}
              placeholder="e.g., Lecture 5: Introduction to Algorithms"
            />
          </div>

          <div>
            <Label htmlFor="description">Description (optional)</Label>
            <Textarea
              id="description"
              {...register('description')}
              placeholder="What will be covered in this meeting?"
            />
          </div>

          <div className="grid grid-cols-2 gap-4">
            <div>
              <Label htmlFor="scheduledStart">Date & Time</Label>
              <Input
                id="scheduledStart"
                type="datetime-local"
                {...register('scheduledStart', { required: true })}
              />
            </div>
            <div>
              <Label htmlFor="durationMinutes">Duration (minutes)</Label>
              <Input
                id="durationMinutes"
                type="number"
                {...register('durationMinutes', { min: 15, max: 480 })}
              />
            </div>
          </div>

          <div className="space-y-3 border-t pt-4">
            <div className="flex items-center justify-between">
              <Label htmlFor="waitingRoom">Waiting Room</Label>
              <Switch id="waitingRoom" {...register('waitingRoom')} />
            </div>
            
            <div className="flex items-center justify-between">
              <Label htmlFor="muteOnEntry">Mute participants on entry</Label>
              <Switch id="muteOnEntry" {...register('muteOnEntry')} />
            </div>
            
            <div className="flex items-center justify-between">
              <Label htmlFor="autoRecording">Auto-record to cloud</Label>
              <Switch id="autoRecording" {...register('autoRecording')} />
            </div>
          </div>

          <div className="flex justify-end gap-2">
            <Button type="button" variant="outline" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit" disabled={loading}>
              {loading ? 'Creating...' : 'Create Meeting'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
```

```typescript
// frontend/src/components/VideoMeeting/RecordingsList.tsx
import React from 'react';
import { Play, Download, Clock, FileVideo } from 'lucide-react';
import { formatDuration, formatDate } from '@/lib/utils';

interface Recording {
  id: string;
  recordingType: 'video' | 'audio' | 'transcript';
  fileType: string;
  durationSeconds: number;
  fileSizeBytes: number;
  status: 'processing' | 'ready' | 'failed';
  createdAt: string;
}

interface RecordingsListProps {
  meetingId: string;
  recordings: Recording[];
}

export function RecordingsList({ meetingId, recordings }: RecordingsListProps) {
  const handlePlay = async (recordingId: string) => {
    const response = await fetch(`/api/v1/meetings/${meetingId}/recordings/${recordingId}/play`);
    const { playUrl } = await response.json();
    window.open(playUrl, '_blank');
  };

  const handleDownload = async (recordingId: string) => {
    const response = await fetch(`/api/v1/meetings/${meetingId}/recordings/${recordingId}/download`);
    const { downloadUrl } = await response.json();
    window.location.href = downloadUrl;
  };

  if (recordings.length === 0) {
    return (
      <div className="text-center py-8 text-gray-500">
        No recordings available for this meeting.
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {recordings.map((recording) => (
        <div
          key={recording.id}
          className="flex items-center justify-between p-4 border rounded-lg"
        >
          <div className="flex items-center gap-3">
            <FileVideo className="w-8 h-8 text-blue-600" />
            <div>
              <div className="font-medium">
                {recording.recordingType === 'video' ? 'Video Recording' : 
                 recording.recordingType === 'audio' ? 'Audio Only' : 'Transcript'}
              </div>
              <div className="text-sm text-gray-500 flex items-center gap-2">
                <Clock className="w-3 h-3" />
                {formatDuration(recording.durationSeconds)}
                <span className="mx-1">•</span>
                {formatFileSize(recording.fileSizeBytes)}
                <span className="mx-1">•</span>
                {formatDate(recording.createdAt)}
              </div>
            </div>
          </div>

          <div className="flex gap-2">
            {recording.status === 'ready' ? (
              <>
                <button
                  onClick={() => handlePlay(recording.id)}
                  className="p-2 rounded-full hover:bg-gray-100"
                  title="Play"
                >
                  <Play className="w-5 h-5" />
                </button>
                <button
                  onClick={() => handleDownload(recording.id)}
                  className="p-2 rounded-full hover:bg-gray-100"
                  title="Download"
                >
                  <Download className="w-5 h-5" />
                </button>
              </>
            ) : recording.status === 'processing' ? (
              <span className="text-sm text-yellow-600">Processing...</span>
            ) : (
              <span className="text-sm text-red-600">Failed</span>
            )}
          </div>
        </div>
      ))}
    </div>
  );
}

function formatFileSize(bytes: number): string {
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  if (bytes < 1024 * 1024 * 1024) return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
  return `${(bytes / (1024 * 1024 * 1024)).toFixed(2)} GB`;
}
```

---

### 5.6 Сложность освоения и реализации

#### Оценка сложности по провайдерам:

| Провайдер | API сложность | Документация | Настройка | Общая |
|-----------|---------------|--------------|-----------|-------|
| **Zoom** | 🟡 Средняя | 🟢 Отличная | 🟡 OAuth | 🟡 Средняя |
| **BigBlueButton** | 🟢 Простая | 🟡 Средняя | 🔴 Сервер | 🟡 Средняя |
| **Jitsi** | 🟢 Простая | 🟢 Хорошая | 🟢 Просто | 🟢 Низкая |
| **Teams** | 🔴 Сложная | 🟡 Средняя | 🔴 Permissions | 🔴 Высокая |
| **Google Meet** | 🔴 Сложная | 🟡 Средняя | 🔴 Workspace | 🔴 Высокая |

#### Временные оценки:

```
Zoom Integration (рекомендуемый старт):
├── OAuth + API setup: 2-3 дня
├── Create/manage meetings: 3-4 дня
├── Join URL generation: 1-2 дня
├── Webhooks (attendance): 2-3 дня
├── Recordings API: 2-3 дня
├── Frontend components: 3-4 дня
├── Testing: 2-3 дня
└── Итого: 2.5-3 недели

BigBlueButton (open source):
├── Server setup (Docker): 1-2 дня
├── API integration: 2-3 дня
├── Join/create: 2 дня
├── Recordings: 2-3 дня
├── Frontend: 2-3 дня
├── Testing: 2 дня
└── Итого: 2 недели

Полная система (Zoom + BBB + UI):
└── Итого: 5-6 недель
```

#### Типичные проблемы:

| Проблема | Причина | Решение |
|----------|---------|---------|
| OAuth tokens expire | 1 час lifetime | Refresh token flow |
| Webhook delivery | Firewall, HTTPS | Webhook relay service |
| Recording delay | Processing time | Status polling + webhooks |
| Timezone issues | UTC vs local | Всегда UTC на backend |
| Rate limits | API throttling | Queue + backoff |

---

### 5.7 Источники для дальнейшего изучения

#### Официальная документация:

| Провайдер | Ссылка | Ключевые разделы |
|-----------|--------|------------------|
| **Zoom** | [developers.zoom.us](https://developers.zoom.us/docs/) | OAuth, Meetings API, Webhooks |
| **BigBlueButton** | [docs.bigbluebutton.org](https://docs.bigbluebutton.org/) | API, Recording, Greenlight |
| **Jitsi** | [jitsi.github.io](https://jitsi.github.io/handbook/) | Self-hosting, JWT auth |
| **Teams** | [docs.microsoft.com](https://docs.microsoft.com/en-us/graph/api/resources/onlinemeeting) | Graph API, Permissions |
| **Daily.co** | [docs.daily.co](https://docs.daily.co/) | REST API, Prebuilt UI |

#### SDK и библиотеки:

| Язык | Библиотека | Ссылка |
|------|------------|--------|
| **Go** | zoom-lib-golang | [github.com/himalayan-institute](https://github.com/himalayan-institute/zoom-lib-golang) |
| **Go** | bigbluebutton-api-go | [github.com/blindsidenetworks](https://github.com/blindsidenetworks/bigbluebutton-api-go) |
| **JS** | @zoom/meetingsdk | [npmjs.com](https://www.npmjs.com/package/@zoom/meetingsdk) |
| **JS** | jitsi-meet-react | [npmjs.com](https://www.npmjs.com/package/@jitsi/react-sdk) |

#### Примеры реализаций:

| Проект | Описание | Ссылка |
|--------|----------|--------|
| **Canvas LMS** | Open source LMS с Zoom | [github.com/instructure](https://github.com/instructure/canvas-lms) |
| **Moodle** | BBB интеграция | [moodle.org/plugins](https://moodle.org/plugins/mod_bigbluebuttonbn) |
| **Greenlight** | UI для BBB | [github.com/bigbluebutton](https://github.com/bigbluebutton/greenlight) |

---

### 5.8 Чек-лист реализации

```
Phase 1: Infrastructure (Day 1-5)
□ Database schema (meetings, participants, recordings)
□ Video provider interface/abstraction
□ Configuration per tenant
□ API routes structure

Phase 2: Zoom Integration (Day 6-12)
□ OAuth Server-to-Server setup
□ Create meeting API
□ Get meeting info API
□ Update/delete meeting API
□ Generate join URL with SSO
□ Webhook endpoint
□ Webhook event processing

Phase 3: BigBlueButton (Day 13-18)
□ Docker-compose setup
□ API client implementation
□ Create/join/end meeting
□ Recording retrieval
□ Webhook/callback handling

Phase 4: Frontend (Day 19-25)
□ Meeting button component
□ Create meeting dialog
□ Meeting list in course
□ Upcoming meetings widget
□ Recordings list
□ Attendance view

Phase 5: Integration (Day 26-30)
□ Calendar integration
□ Course schedule sync
□ Attendance → gradebook
□ Notifications (before meeting)
□ Recording auto-publish

Phase 6: Polish (Day 31-35)
□ Mobile responsiveness
□ Error handling
□ Loading states
□ Admin settings page
□ Documentation
□ Testing with real Zoom/BBB
```

---

## 6. LTI 1.3 (Learning Tools Interoperability)

### 6.1 Определение

**LTI (Learning Tools Interoperability)** — это стандарт IMS Global, позволяющий безопасно интегрировать внешние образовательные инструменты (tools) в системы управления обучением (LMS). LTI создает "plug-and-play" экосистему, где любой LTI-совместимый инструмент может работать в любой LTI-совместимой платформе.

**Ключевая идея:** Единый стандарт интеграции вместо тысяч кастомных API. Инструмент реализует LTI один раз и работает в Canvas, Moodle, Blackboard, D2L и любой другой LMS.

#### Эволюция LTI:

| Версия | Год | Особенности | Статус |
|--------|-----|-------------|--------|
| **LTI 1.0** | 2010 | Basic Launch (OAuth 1.0) | Устарела |
| **LTI 1.1** | 2012 | Outcomes (оценки назад) | Широко используется |
| **LTI 1.3** | 2019 | OAuth 2.0, JWT, Services | Текущий стандарт ✅ |
| **LTI 2.0** | 2014 | Сложный, избыточный | Отменен |
| **LTI Advantage** | 2019 | LTI 1.3 + Services | Рекомендуется |

#### Архитектура LTI 1.3:

```
┌─────────────────────────────────────────────────────────────────┐
│                      LTI 1.3 Architecture                        │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                    Platform (LMS)                        │    │
│  │                    [Our Application]                     │    │
│  │                                                          │    │
│  │  • Registers external tools                              │    │
│  │  • Initiates OIDC login flow                            │    │
│  │  • Issues JWT id_token                                   │    │
│  │  • Receives grades via AGS                               │    │
│  │  • Provides roster via NRPS                              │    │
│  │                                                          │    │
│  └──────────────────────────┬───────────────────────────────┘    │
│                             │                                    │
│                   OIDC + JWT Launch                              │
│                             │                                    │
│  ┌──────────────────────────▼───────────────────────────────┐    │
│  │                    Tool Provider                          │    │
│  │              [External Application]                       │    │
│  │                                                          │    │
│  │  Examples:                                                │    │
│  │  • Turnitin (plagiarism)                                 │    │
│  │  • Kahoot (quizzes)                                      │    │
│  │  • Labster (virtual labs)                                │    │
│  │  • H5P (interactive content)                             │    │
│  │  • Proctorio (proctoring)                                │    │
│  │  • LinkedIn Learning                                      │    │
│  │  • McGraw-Hill Connect                                    │    │
│  │                                                          │    │
│  └──────────────────────────────────────────────────────────┘    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

#### Роли в LTI:

| Роль | Описание | Примеры | Наша роль |
|------|----------|---------|-----------|
| **Platform** | LMS, которая запускает tools | Canvas, Moodle, наша LMS | ✅ Platform |
| **Tool** | Внешнее приложение | Turnitin, Kahoot, H5P | ⚠️ Можно быть Tool тоже |

#### Ключевые компоненты LTI 1.3:

| Компонент | Описание | Протокол |
|-----------|----------|----------|
| **OIDC Launch** | Безопасный запуск инструмента | OpenID Connect |
| **JWT id_token** | Данные о пользователе и контексте | JSON Web Token |
| **Deep Linking** | Выбор контента из tool | LTI DL 2.0 |
| **AGS** | Assignment and Grade Services | REST API |
| **NRPS** | Names and Roles Provisioning | REST API |

#### LTI 1.3 Launch Flow:

```
┌────────────────────────────────────────────────────────────────────────┐
│                        LTI 1.3 Launch Flow                              │
├────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│  User clicks "Launch Tool" in LMS                                       │
│           │                                                             │
│           ▼                                                             │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ Step 1: OIDC Login Initiation                                    │   │
│  │                                                                   │   │
│  │ Platform → Tool (via browser redirect)                           │   │
│  │                                                                   │   │
│  │ GET /oidc/login?                                                  │   │
│  │   iss=https://lms.example.com                                    │   │
│  │   &login_hint=user_12345                                         │   │
│  │   &target_link_uri=https://tool.com/launch                       │   │
│  │   &lti_message_hint=context_xyz                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│           │                                                             │
│           ▼                                                             │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ Step 2: Tool Authentication Request                              │   │
│  │                                                                   │   │
│  │ Tool → Platform (via browser redirect)                           │   │
│  │                                                                   │   │
│  │ GET /oauth2/authorize?                                            │   │
│  │   response_type=id_token                                         │   │
│  │   &client_id=tool_client_id                                      │   │
│  │   &redirect_uri=https://tool.com/launch                          │   │
│  │   &scope=openid                                                  │   │
│  │   &state=random_state                                            │   │
│  │   &nonce=random_nonce                                            │   │
│  │   &login_hint=user_12345                                         │   │
│  │   &lti_message_hint=context_xyz                                  │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│           │                                                             │
│           ▼                                                             │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ Step 3: Platform Issues JWT id_token                             │   │
│  │                                                                   │   │
│  │ Platform → Tool (via browser POST)                               │   │
│  │                                                                   │   │
│  │ POST /launch                                                      │   │
│  │ id_token=eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...                 │   │
│  │ state=random_state                                               │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│           │                                                             │
│           ▼                                                             │
│  ┌─────────────────────────────────────────────────────────────────┐   │
│  │ Step 4: Tool Validates JWT & Renders Content                     │   │
│  │                                                                   │   │
│  │ Tool:                                                             │   │
│  │   1. Validates JWT signature (using platform's public key)       │   │
│  │   2. Validates claims (iss, aud, nonce, exp)                     │   │
│  │   3. Extracts user info, course, role                            │   │
│  │   4. Renders appropriate content                                 │   │
│  └─────────────────────────────────────────────────────────────────┘   │
│                                                                         │
└────────────────────────────────────────────────────────────────────────┘
```

#### Структура JWT id_token (LTI 1.3):

```json
{
  "iss": "https://lms.example.com",
  "sub": "user_12345",
  "aud": ["tool_client_id"],
  "exp": 1609459200,
  "iat": 1609455600,
  "nonce": "abc123",
  
  "https://purl.imsglobal.org/spec/lti/claim/message_type": "LtiResourceLinkRequest",
  "https://purl.imsglobal.org/spec/lti/claim/version": "1.3.0",
  
  "https://purl.imsglobal.org/spec/lti/claim/deployment_id": "deployment_1",
  
  "https://purl.imsglobal.org/spec/lti/claim/target_link_uri": "https://tool.com/content/123",
  
  "https://purl.imsglobal.org/spec/lti/claim/resource_link": {
    "id": "resource_link_123",
    "title": "Chapter 5 Quiz",
    "description": "Quiz on algorithms"
  },
  
  "https://purl.imsglobal.org/spec/lti/claim/roles": [
    "http://purl.imsglobal.org/vocab/lis/v2/membership#Learner",
    "http://purl.imsglobal.org/vocab/lis/v2/institution/person#Student"
  ],
  
  "https://purl.imsglobal.org/spec/lti/claim/context": {
    "id": "course_456",
    "label": "CS101",
    "title": "Introduction to Computer Science",
    "type": ["http://purl.imsglobal.org/vocab/lis/v2/course#CourseOffering"]
  },
  
  "https://purl.imsglobal.org/spec/lti/claim/lis": {
    "person_sourcedid": "student_12345",
    "course_section_sourcedid": "CS101-001"
  },
  
  "name": "Иван Иванов",
  "given_name": "Иван",
  "family_name": "Иванов",
  "email": "ivan@example.com",
  "picture": "https://lms.example.com/avatars/12345.jpg",
  
  "https://purl.imsglobal.org/spec/lti-ags/claim/endpoint": {
    "scope": [
      "https://purl.imsglobal.org/spec/lti-ags/scope/lineitem",
      "https://purl.imsglobal.org/spec/lti-ags/scope/score"
    ],
    "lineitems": "https://lms.example.com/api/lti/courses/456/lineitems",
    "lineitem": "https://lms.example.com/api/lti/courses/456/lineitems/789"
  },
  
  "https://purl.imsglobal.org/spec/lti-nrps/claim/namesroleservice": {
    "context_memberships_url": "https://lms.example.com/api/lti/courses/456/memberships",
    "service_versions": ["2.0"]
  }
}
```

#### LTI Advantage Services:

| Service | Аббревиатура | Описание | Use Case |
|---------|--------------|----------|----------|
| **Assignment and Grade Services** | AGS | Отправка оценок в LMS | Tool → Platform: "Студент получил 85/100" |
| **Names and Role Provisioning** | NRPS | Получение списка студентов | Tool запрашивает roster курса |
| **Deep Linking** | DL 2.0 | Выбор контента из tool | Преподаватель выбирает quiz из Kahoot |

---

### 6.2 Почему важно реализовать

#### Бизнес-причины:

| Причина | Описание | Влияние |
|---------|----------|---------|
| **Экосистема 5000+ tools** | Огромный выбор интеграций | Instant value для клиентов |
| **Индустриальный стандарт** | Все LMS конкуренты поддерживают | Без LTI = неконкурентоспособны |
| **Сертификация IMS** | Знак качества | Доверие enterprise клиентов |
| **EdTech partnerships** | Партнерства с вендорами | Совместный маркетинг |
| **Vendor independence** | Клиенты не зависят от наших tools | Свобода выбора |

#### Технические причины:

```
📊 Статистика LTI (2024):

• 5000+ LTI-совместимых инструментов
• 99% LMS систем поддерживают LTI
• $50B+ рынок EdTech tools
• Top tools по популярности:
  - Turnitin (90% университетов)
  - Kahoot (9M+ учителей)
  - H5P (open source интерактивный контент)
  - Labster (virtual labs)
  - Proctorio/Respondus (proctoring)
  - Publisher content (Pearson, McGraw-Hill)

💡 Почему Tool Provider хотят LTI:
• Один API для всех LMS (не 20 кастомных)
• SSO из коробки (нет отдельного логина)
• Контекст (курс, роль) автоматически
• Оценки обратно в LMS (AGS)
```

#### Сравнение: с LTI vs без LTI

```
Без LTI (кастомная интеграция):
├─ Каждый tool = отдельный договор
├─ Каждый tool = кастомный код интеграции
├─ Каждый tool = отдельный логин для студентов
├─ Оценки не синхронизируются (копировать вручную)
├─ 10 tools = 10 интеграций = 10 месяцев работы
└─ Обновление tool = переделка интеграции

С LTI:
├─ Один стандарт = любой LTI tool работает
├─ Добавить tool = 5 минут конфигурации
├─ SSO автоматический
├─ Оценки автоматически в журнале (AGS)
├─ 10 tools = 10 × 5 минут = 50 минут ✅
└─ Обновление tool = прозрачно для нас
```

#### Примеры реальных интеграций:

```
Сценарий 1: Университет с Turnitin
├─ Без LTI: 
│   • Студенты загружают в Turnitin отдельно
│   • Результаты копируют вручную
│   • 30 мин/работу × 100 работ = 50 часов
├─ С LTI:
│   • Студент сдает через LMS
│   • Проверка plagiarism автоматически
│   • Оценка автоматически в журнале
│   • 50 часов → 0 часов ручной работы ✅

Сценарий 2: Корпорация с LinkedIn Learning
├─ Без LTI:
│   • Отдельный логин в LinkedIn Learning
│   • Прохождение не отслеживается
│   • ROI неизвестен
├─ С LTI:
│   • SSO из корпоративной LMS
│   • Completion автоматически в LMS
│   • Единый отчет по обучению ✅
```

---

### 6.3 Что дает конечному пользователю

#### Для студентов:

| Функция | Описание | Польза |
|---------|----------|--------|
| **Single Sign-On** | Один логин для всего | Нет множества паролей |
| **Единый интерфейс** | Tools внутри LMS | Не нужно переключаться |
| **Автоматические оценки** | Результаты в журнале | Видят прогресс сразу |
| **Богатый контент** | Интерактивные tools | Engaging обучение |
| **Мобильный доступ** | Tools через LMS app | Удобство |

#### Для преподавателей:

| Функция | Описание | Польза |
|---------|----------|--------|
| **Easy setup** | Добавить tool = 5 минут | Экономия времени |
| **Deep Linking** | Выбрать конкретный контент | Гибкость |
| **Grade sync** | Оценки автоматически | Нет ручного ввода |
| **Roster sync** | Студенты автоматически | Актуальные списки |
| **Analytics** | Данные о использовании | Insights |

#### Для администраторов:

| Функция | Описание | Польза |
|---------|----------|--------|
| **Centralized management** | Все tools в одном месте | Контроль |
| **Security** | OAuth 2.0, JWT, HTTPS | Compliance |
| **Audit trail** | Логи запусков | Отчетность |
| **Cost control** | Видимость использования | ROI |

#### Популярные LTI Tools по категориям:

```
📝 Assessment & Proctoring:
• Turnitin - plagiarism detection
• Proctorio - online proctoring
• Respondus LockDown Browser
• Examity - remote proctoring
• Gradescope - AI grading

🎮 Interactive Content:
• H5P - interactive exercises
• Kahoot - gamified quizzes
• Nearpod - interactive lessons
• Pear Deck - engagement
• Edpuzzle - interactive video

🔬 STEM & Labs:
• Labster - virtual labs
• PhET Simulations
• MATLAB Grader
• WebAssign - homework
• Cengage WebAssign

📚 Publisher Content:
• Pearson MyLab
• McGraw-Hill Connect
• Cengage MindTap
• Wiley Plus
• Macmillan Achieve

💼 Professional Development:
• LinkedIn Learning
• Coursera for Business
• Udemy Business
• Pluralsight
• Skillsoft

🎥 Video & Collaboration:
• Kaltura - video platform
• Panopto - lecture capture
• VoiceThread - discussions
• Padlet - collaboration
• Flipgrid - video discussions
```

---

### 6.4 Какие пользователи предпочитают эту функцию

#### Целевые сегменты:

| Сегмент | Важность LTI | Типичные tools |
|---------|--------------|----------------|
| **Университеты** | Критическая | Turnitin, Labster, Publisher |
| **K-12** | Высокая | Kahoot, Nearpod, Edpuzzle |
| **Корпорации** | Высокая | LinkedIn Learning, Skillsoft |
| **Bootcamps** | Средняя | Собственные tools |
| **Репетиторство** | Низкая | Редко нужно |

#### Типичные требования по сегментам:

```
🎓 Университеты:
• "Нужна интеграция с Turnitin" (99% случаев)
• "Используем Pearson MyLab для учебников"
• "Нужен Respondus для экзаменов"
• "Labster для естественных наук"
• IMS сертификация = обязательно для тендеров

💼 Корпорации:
• "Интегрируйтесь с LinkedIn Learning"
• "Нужен Skillsoft каталог"
• "Content от провайдеров compliance training"
• Security review = обязательно

🏫 K-12:
• "Учителя уже используют Kahoot"
• "Нужен Nearpod для интерактивных уроков"
• "Google Classroom совместимость"
• Простота = критично
```

#### Вопросы при продаже:

```
Типичные вопросы от клиентов:

1. "Поддерживаете ли вы LTI 1.3?" (must have)
2. "Есть ли LTI Advantage (AGS, NRPS, Deep Linking)?"
3. "Можно ли добавить Turnitin?" (90% университетов)
4. "Работает ли интеграция с publisher content?"
5. "Есть ли IMS сертификация?"
6. "Можем ли мы быть Tool Provider тоже?"
7. "Как работает grade passback?"
8. "Поддерживается ли Deep Linking?"

Без LTI 1.3:
• Потеря 95%+ образовательных клиентов
• Невозможность участвовать в тендерах
• Нет интеграции с Turnitin = deal breaker
```

---

### 6.5 Как интегрировать в наше приложение

#### Два режима работы:

```
┌─────────────────────────────────────────────────────────────────┐
│                    LTI Integration Modes                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Mode 1: Platform (Primary)                                      │
│  ───────────────────────────                                     │
│  Our LMS launches external tools                                 │
│  • Students use Turnitin from our LMS                           │
│  • We receive grades from tools (AGS)                           │
│  • We provide roster to tools (NRPS)                            │
│                                                                  │
│  Mode 2: Tool Provider (Secondary)                               │
│  ─────────────────────────────────                               │
│  Our content/features available in other LMS                    │
│  • Our quiz engine in Canvas/Moodle                             │
│  • Our content library in partner LMS                           │
│  • Enterprise customers embed our features                       │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

#### Архитектура Platform:

```
┌─────────────────────────────────────────────────────────────────┐
│                  LTI Platform Architecture                       │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                   Tool Registration                      │    │
│  │                                                          │    │
│  │  POST /api/admin/lti/tools                               │    │
│  │  {                                                       │    │
│  │    "name": "Turnitin",                                   │    │
│  │    "client_id": "turnitin_abc123",                       │    │
│  │    "oidc_login_url": "https://turnitin.com/lti/login",  │    │
│  │    "target_link_uri": "https://turnitin.com/lti/launch",│    │
│  │    "jwks_url": "https://turnitin.com/.well-known/jwks", │    │
│  │    "scopes": ["AGS", "NRPS", "DeepLinking"]             │    │
│  │  }                                                       │    │
│  └─────────────────────────────────────────────────────────┘    │
│                             │                                    │
│                             ▼                                    │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                   Launch Endpoints                       │    │
│  │                                                          │    │
│  │  GET  /lti/oidc/login          ← Tool redirects here    │    │
│  │  POST /lti/oidc/authorize      ← We issue id_token      │    │
│  │  GET  /lti/jwks                ← Our public keys        │    │
│  │                                                          │    │
│  └─────────────────────────────────────────────────────────┘    │
│                             │                                    │
│                             ▼                                    │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                   LTI Advantage Services                 │    │
│  │                                                          │    │
│  │  AGS (Assignment and Grade Services):                    │    │
│  │  GET  /lti/courses/:id/lineitems                        │    │
│  │  POST /lti/courses/:id/lineitems                        │    │
│  │  POST /lti/courses/:id/lineitems/:id/scores             │    │
│  │                                                          │    │
│  │  NRPS (Names and Role Provisioning):                    │    │
│  │  GET  /lti/courses/:id/memberships                      │    │
│  │                                                          │    │
│  │  Deep Linking:                                           │    │
│  │  POST /lti/deep-linking/callback                        │    │
│  │                                                          │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

#### Схема базы данных:

```sql
-- Зарегистрированные LTI Tools
CREATE TABLE lti_tools (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    
    -- Идентификация
    name VARCHAR(255) NOT NULL,
    description TEXT,
    icon_url VARCHAR(500),
    
    -- LTI 1.3 Configuration
    client_id VARCHAR(255) NOT NULL,          -- Наш client_id для этого tool
    deployment_id VARCHAR(255) NOT NULL,      -- Уникальный для каждого deployment
    
    -- Tool URLs
    oidc_login_url VARCHAR(500) NOT NULL,     -- Куда отправлять OIDC login
    target_link_uri VARCHAR(500) NOT NULL,    -- Default launch URL
    redirect_uris TEXT[],                      -- Allowed redirect URIs
    
    -- Tool's public keys
    jwks_url VARCHAR(500),                    -- URL для получения public keys
    public_key TEXT,                          -- Или статический public key
    
    -- Capabilities
    supports_deep_linking BOOLEAN DEFAULT false,
    supports_ags BOOLEAN DEFAULT false,       -- Assignment and Grade Services
    supports_nrps BOOLEAN DEFAULT false,      -- Names and Role Provisioning
    
    -- Scopes granted
    scopes TEXT[] DEFAULT '{}',
    /*
    Possible scopes:
    - https://purl.imsglobal.org/spec/lti-ags/scope/lineitem
    - https://purl.imsglobal.org/spec/lti-ags/scope/lineitem.readonly
    - https://purl.imsglobal.org/spec/lti-ags/scope/result.readonly
    - https://purl.imsglobal.org/spec/lti-ags/scope/score
    - https://purl.imsglobal.org/spec/lti-nrps/scope/contextmembership.readonly
    */
    
    -- Custom parameters
    custom_parameters JSONB DEFAULT '{}',
    
    -- Privacy settings
    send_name BOOLEAN DEFAULT true,
    send_email BOOLEAN DEFAULT true,
    send_avatar BOOLEAN DEFAULT false,
    
    -- Status
    is_enabled BOOLEAN DEFAULT true,
    is_global BOOLEAN DEFAULT false,          -- Available for all courses
    
    -- Metadata
    created_by UUID REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(tenant_id, client_id)
);

-- Tool placement in courses
CREATE TABLE lti_tool_placements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tool_id UUID NOT NULL REFERENCES lti_tools(id) ON DELETE CASCADE,
    
    -- Context
    course_id UUID NOT NULL REFERENCES courses(id) ON DELETE CASCADE,
    module_id UUID,                           -- Опционально, если в конкретном модуле
    
    -- Launch configuration
    resource_link_id VARCHAR(255) NOT NULL,   -- Уникальный ID этого размещения
    title VARCHAR(255) NOT NULL,
    description TEXT,
    
    -- Custom launch URL (если отличается от default)
    custom_launch_url VARCHAR(500),
    custom_parameters JSONB DEFAULT '{}',
    
    -- Display settings
    display_type VARCHAR(50) DEFAULT 'iframe', -- 'iframe', 'new_window', 'embed'
    iframe_width INTEGER,
    iframe_height INTEGER,
    
    -- Deep Linking content (если выбран через DL)
    deep_link_content JSONB,
    
    -- Grade settings
    lineitem_id UUID,                         -- Связь с gradebook
    max_score DECIMAL(10,2),
    
    -- Order
    sort_order INTEGER DEFAULT 0,
    
    -- Status
    is_visible BOOLEAN DEFAULT true,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(course_id, resource_link_id)
);

-- LTI Launch sessions (для отслеживания)
CREATE TABLE lti_launches (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Context
    tool_id UUID NOT NULL REFERENCES lti_tools(id),
    placement_id UUID REFERENCES lti_tool_placements(id),
    user_id UUID NOT NULL REFERENCES users(id),
    course_id UUID REFERENCES courses(id),
    
    -- Launch data
    message_type VARCHAR(100) NOT NULL,       -- 'LtiResourceLinkRequest', 'LtiDeepLinkingRequest'
    nonce VARCHAR(255) NOT NULL,              -- Для предотвращения replay
    state VARCHAR(255),
    
    -- Status
    status VARCHAR(50) DEFAULT 'initiated',   -- 'initiated', 'completed', 'failed'
    
    -- Timestamps
    initiated_at TIMESTAMP DEFAULT NOW(),
    completed_at TIMESTAMP,
    
    -- Error tracking
    error_message TEXT,
    
    -- Expiration (launches expire after 10 minutes)
    expires_at TIMESTAMP DEFAULT NOW() + INTERVAL '10 minutes'
);

-- AGS Line Items (связь с gradebook)
CREATE TABLE lti_lineitems (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Context
    tool_id UUID NOT NULL REFERENCES lti_tools(id),
    course_id UUID NOT NULL REFERENCES courses(id),
    placement_id UUID REFERENCES lti_tool_placements(id),
    
    -- Line item data
    label VARCHAR(255) NOT NULL,
    score_maximum DECIMAL(10,2) NOT NULL,
    tag VARCHAR(255),                         -- Optional categorization
    
    -- Resource link (если привязан к конкретному placement)
    resource_link_id VARCHAR(255),
    
    -- Связь с нашим gradebook
    gradebook_column_id UUID,                 -- FK to our grade columns
    
    -- Timestamps
    start_date TIMESTAMP,
    end_date TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- AGS Scores (оценки от tools)
CREATE TABLE lti_scores (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    lineitem_id UUID NOT NULL REFERENCES lti_lineitems(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id),
    
    -- Score data
    score_given DECIMAL(10,2),
    score_maximum DECIMAL(10,2),
    
    -- Activity progress
    activity_progress VARCHAR(50),            -- 'Initialized', 'Started', 'InProgress', 'Submitted', 'Completed'
    grading_progress VARCHAR(50),             -- 'FullyGraded', 'Pending', 'PendingManual', 'Failed', 'NotReady'
    
    -- Comment
    comment TEXT,
    
    -- Timestamps
    timestamp TIMESTAMP DEFAULT NOW(),
    
    -- Для синхронизации с нашим gradebook
    synced_to_gradebook BOOLEAN DEFAULT false,
    synced_at TIMESTAMP,
    
    UNIQUE(lineitem_id, user_id)
);

-- Platform keys (RSA keys for signing JWTs)
CREATE TABLE lti_platform_keys (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    
    -- Key pair
    kid VARCHAR(255) NOT NULL,                -- Key ID
    algorithm VARCHAR(10) DEFAULT 'RS256',
    public_key TEXT NOT NULL,
    private_key TEXT NOT NULL,                -- Encrypted
    
    -- Status
    is_active BOOLEAN DEFAULT true,
    
    -- Rotation
    created_at TIMESTAMP DEFAULT NOW(),
    expires_at TIMESTAMP,
    
    UNIQUE(tenant_id, kid)
);

-- Indexes
CREATE INDEX idx_lti_tools_tenant ON lti_tools(tenant_id);
CREATE INDEX idx_lti_placements_course ON lti_tool_placements(course_id);
CREATE INDEX idx_lti_placements_tool ON lti_tool_placements(tool_id);
CREATE INDEX idx_lti_launches_user ON lti_launches(user_id);
CREATE INDEX idx_lti_launches_nonce ON lti_launches(nonce);
CREATE INDEX idx_lti_lineitems_course ON lti_lineitems(course_id);
CREATE INDEX idx_lti_scores_lineitem ON lti_scores(lineitem_id);
CREATE INDEX idx_lti_scores_user ON lti_scores(user_id);
```

#### Frontend компоненты:

```typescript
// frontend/src/components/LTI/LaunchButton.tsx
import React, { useState } from 'react';
import { ExternalLink, Loader2 } from 'lucide-react';
import { Button } from '@/components/ui/button';

interface LaunchButtonProps {
  placementId: string;
  title: string;
  displayType: 'iframe' | 'new_window' | 'embed';
}

export function LTILaunchButton({ placementId, title, displayType }: LaunchButtonProps) {
  const [launching, setLaunching] = useState(false);
  const [iframeUrl, setIframeUrl] = useState<string | null>(null);

  const handleLaunch = async () => {
    setLaunching(true);
    try {
      // Получаем launch URL
      const response = await fetch(`/api/v1/lti/placements/${placementId}/launch`, {
        method: 'POST',
      });
      const { launchUrl, method, formData } = await response.json();

      if (displayType === 'new_window') {
        // Открываем в новом окне
        if (method === 'POST') {
          // Создаем скрытую форму для POST
          const form = document.createElement('form');
          form.method = 'POST';
          form.action = launchUrl;
          form.target = '_blank';
          
          Object.entries(formData).forEach(([key, value]) => {
            const input = document.createElement('input');
            input.type = 'hidden';
            input.name = key;
            input.value = value as string;
            form.appendChild(input);
          });
          
          document.body.appendChild(form);
          form.submit();
          document.body.removeChild(form);
        } else {
          window.open(launchUrl, '_blank');
        }
      } else {
        // Показываем в iframe
        setIframeUrl(launchUrl);
      }
    } catch (error) {
      console.error('Launch failed:', error);
    } finally {
      setLaunching(false);
    }
  };

  if (iframeUrl) {
    return (
      <div className="w-full h-[600px] border rounded-lg overflow-hidden">
        <div className="bg-gray-100 p-2 flex justify-between items-center">
          <span className="font-medium">{title}</span>
          <Button variant="ghost" size="sm" onClick={() => setIframeUrl(null)}>
            Close
          </Button>
        </div>
        <iframe
          src={iframeUrl}
          className="w-full h-[calc(100%-40px)]"
          sandbox="allow-scripts allow-same-origin allow-forms allow-popups"
          title={title}
        />
      </div>
    );
  }

  return (
    <Button onClick={handleLaunch} disabled={launching}>
      {launching ? (
        <Loader2 className="w-4 h-4 mr-2 animate-spin" />
      ) : (
        <ExternalLink className="w-4 h-4 mr-2" />
      )}
      {title}
    </Button>
  );
}
```

```typescript
// frontend/src/components/LTI/AddToolDialog.tsx
import React, { useState } from 'react';
import { useForm } from 'react-hook-form';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Switch } from '@/components/ui/switch';
import { Label } from '@/components/ui/label';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';

interface AddToolForm {
  name: string;
  description?: string;
  clientId: string;
  oidcLoginUrl: string;
  targetLinkUri: string;
  jwksUrl?: string;
  publicKey?: string;
  supportsDeepLinking: boolean;
  supportsAgs: boolean;
  supportsNrps: boolean;
  sendName: boolean;
  sendEmail: boolean;
}

interface AddToolDialogProps {
  open: boolean;
  onClose: () => void;
  onCreated: (tool: any) => void;
}

export function AddToolDialog({ open, onClose, onCreated }: AddToolDialogProps) {
  const [loading, setLoading] = useState(false);
  const [keyType, setKeyType] = useState<'jwks' | 'public'>('jwks');
  
  const { register, handleSubmit, watch, formState: { errors } } = useForm<AddToolForm>({
    defaultValues: {
      supportsDeepLinking: true,
      supportsAgs: true,
      supportsNrps: false,
      sendName: true,
      sendEmail: true,
    },
  });

  const onSubmit = async (data: AddToolForm) => {
    setLoading(true);
    try {
      const response = await fetch('/api/v1/admin/lti/tools', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: data.name,
          description: data.description,
          clientId: data.clientId,
          oidcLoginUrl: data.oidcLoginUrl,
          targetLinkUri: data.targetLinkUri,
          jwksUrl: keyType === 'jwks' ? data.jwksUrl : undefined,
          publicKey: keyType === 'public' ? data.publicKey : undefined,
          supportsDeepLinking: data.supportsDeepLinking,
          supportsAgs: data.supportsAgs,
          supportsNrps: data.supportsNrps,
          sendName: data.sendName,
          sendEmail: data.sendEmail,
        }),
      });
      
      const tool = await response.json();
      onCreated(tool);
      onClose();
    } catch (error) {
      console.error('Failed to add tool:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[600px] max-h-[90vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Add LTI 1.3 Tool</DialogTitle>
        </DialogHeader>
        
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div className="space-y-4">
            <h4 className="font-medium">Basic Information</h4>
            
            <div>
              <Label htmlFor="name">Tool Name *</Label>
              <Input
                id="name"
                {...register('name', { required: 'Name is required' })}
                placeholder="e.g., Turnitin, Kahoot, H5P"
              />
            </div>

            <div>
              <Label htmlFor="description">Description</Label>
              <Textarea
                id="description"
                {...register('description')}
                placeholder="What this tool does..."
              />
            </div>
          </div>

          <div className="space-y-4 border-t pt-4">
            <h4 className="font-medium">LTI Configuration</h4>
            <p className="text-sm text-gray-500">
              Get these values from the tool provider's LTI configuration page
            </p>
            
            <div>
              <Label htmlFor="clientId">Client ID *</Label>
              <Input
                id="clientId"
                {...register('clientId', { required: 'Client ID is required' })}
                placeholder="Provided by tool vendor"
              />
            </div>

            <div>
              <Label htmlFor="oidcLoginUrl">OIDC Login URL *</Label>
              <Input
                id="oidcLoginUrl"
                {...register('oidcLoginUrl', { required: 'OIDC Login URL is required' })}
                placeholder="https://tool.example.com/lti/login"
              />
            </div>

            <div>
              <Label htmlFor="targetLinkUri">Target Link URI *</Label>
              <Input
                id="targetLinkUri"
                {...register('targetLinkUri', { required: 'Target Link URI is required' })}
                placeholder="https://tool.example.com/lti/launch"
              />
            </div>

            <div>
              <Label>Public Key Configuration</Label>
              <Tabs value={keyType} onValueChange={(v) => setKeyType(v as 'jwks' | 'public')}>
                <TabsList className="w-full">
                  <TabsTrigger value="jwks" className="flex-1">JWKS URL</TabsTrigger>
                  <TabsTrigger value="public" className="flex-1">Public Key</TabsTrigger>
                </TabsList>
                <TabsContent value="jwks">
                  <Input
                    {...register('jwksUrl')}
                    placeholder="https://tool.example.com/.well-known/jwks.json"
                  />
                </TabsContent>
                <TabsContent value="public">
                  <Textarea
                    {...register('publicKey')}
                    placeholder="-----BEGIN PUBLIC KEY-----\n..."
                    rows={4}
                  />
                </TabsContent>
              </Tabs>
            </div>
          </div>

          <div className="space-y-4 border-t pt-4">
            <h4 className="font-medium">LTI Advantage Services</h4>
            
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <div>
                  <Label htmlFor="supportsDeepLinking">Deep Linking</Label>
                  <p className="text-xs text-gray-500">Allow instructors to select specific content</p>
                </div>
                <Switch id="supportsDeepLinking" {...register('supportsDeepLinking')} />
              </div>
              
              <div className="flex items-center justify-between">
                <div>
                  <Label htmlFor="supportsAgs">Assignment & Grades (AGS)</Label>
                  <p className="text-xs text-gray-500">Sync grades back to gradebook</p>
                </div>
                <Switch id="supportsAgs" {...register('supportsAgs')} />
              </div>
              
              <div className="flex items-center justify-between">
                <div>
                  <Label htmlFor="supportsNrps">Names & Roles (NRPS)</Label>
                  <p className="text-xs text-gray-500">Share course roster with tool</p>
                </div>
                <Switch id="supportsNrps" {...register('supportsNrps')} />
              </div>
            </div>
          </div>

          <div className="space-y-4 border-t pt-4">
            <h4 className="font-medium">Privacy Settings</h4>
            
            <div className="space-y-3">
              <div className="flex items-center justify-between">
                <Label htmlFor="sendName">Send user name</Label>
                <Switch id="sendName" {...register('sendName')} />
              </div>
              
              <div className="flex items-center justify-between">
                <Label htmlFor="sendEmail">Send user email</Label>
                <Switch id="sendEmail" {...register('sendEmail')} />
              </div>
            </div>
          </div>

          <div className="flex justify-end gap-2 pt-4">
            <Button type="button" variant="outline" onClick={onClose}>
              Cancel
            </Button>
            <Button type="submit" disabled={loading}>
              {loading ? 'Adding...' : 'Add Tool'}
            </Button>
          </div>
        </form>
      </DialogContent>
    </Dialog>
  );
}
```

```typescript
// frontend/src/components/LTI/DeepLinkingPicker.tsx
import React, { useEffect, useRef, useState } from 'react';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '@/components/ui/dialog';

interface DeepLinkingPickerProps {
  toolId: string;
  courseId: string;
  open: boolean;
  onClose: () => void;
  onContentSelected: (content: DeepLinkContent) => void;
}

interface DeepLinkContent {
  type: string;
  title: string;
  url: string;
  custom?: Record<string, string>;
  lineItem?: {
    scoreMaximum: number;
    label: string;
  };
}

export function DeepLinkingPicker({ 
  toolId, 
  courseId, 
  open, 
  onClose, 
  onContentSelected 
}: DeepLinkingPickerProps) {
  const iframeRef = useRef<HTMLIFrameElement>(null);
  const [launchUrl, setLaunchUrl] = useState<string | null>(null);

  useEffect(() => {
    if (open) {
      // Инициируем Deep Linking launch
      fetch(`/api/v1/lti/tools/${toolId}/deep-linking`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ courseId }),
      })
        .then(res => res.json())
        .then(data => setLaunchUrl(data.launchUrl));
    }
  }, [open, toolId, courseId]);

  useEffect(() => {
    // Слушаем postMessage от tool
    const handleMessage = (event: MessageEvent) => {
      if (event.data.type === 'lti_deep_linking_response') {
        onContentSelected(event.data.content);
        onClose();
      }
    };

    window.addEventListener('message', handleMessage);
    return () => window.removeEventListener('message', handleMessage);
  }, [onContentSelected, onClose]);

  return (
    <Dialog open={open} onOpenChange={onClose}>
      <DialogContent className="sm:max-w-[900px] h-[80vh]">
        <DialogHeader>
          <DialogTitle>Select Content</DialogTitle>
        </DialogHeader>
        
        {launchUrl ? (
          <iframe
            ref={iframeRef}
            src={launchUrl}
            className="w-full h-full border-0"
            sandbox="allow-scripts allow-same-origin allow-forms allow-popups"
            title="Deep Linking Picker"
          />
        ) : (
          <div className="flex items-center justify-center h-full">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600" />
          </div>
        )}
      </DialogContent>
    </Dialog>
  );
}
```

---

### 6.6 Сложность освоения и реализации

#### Оценка сложности:

| Аспект | Сложность | Комментарий |
|--------|-----------|-------------|
| **OIDC Flow** | 🔴 Высокая | Много шагов, JWT, security |
| **JWT Generation/Validation** | 🟡 Средняя | Библиотеки помогают |
| **AGS (Grades)** | 🟡 Средняя | REST API, понятная модель |
| **NRPS (Roster)** | 🟢 Низкая | Простой GET endpoint |
| **Deep Linking** | 🟡 Средняя | Дополнительный flow |
| **Security** | 🔴 Высокая | Keys, nonces, validation |
| **Testing** | 🔴 Высокая | Много edge cases |

#### Временные оценки:

```
LTI 1.3 Platform (базовая поддержка):
├── OIDC endpoints: 3-4 дня
├── JWT generation: 2-3 дня
├── JWT validation: 2 дня
├── Tool management: 3-4 дня
├── Launch flow: 3-4 дня
├── Frontend: 3-4 дня
├── Testing: 3-4 дня
└── Итого: 3-4 недели

LTI Advantage:
├── AGS (grades): 3-4 дня
├── NRPS (roster): 2 дня
├── Deep Linking: 3-4 дня
├── Testing: 2-3 дня
└── Итого: +2 недели

IMS Certification preparation:
├── Conformance testing: 1 неделя
├── Bug fixes: 1 неделя
├── Documentation: 3 дня
└── Итого: +2-3 недели

Общее время: 7-9 недель до сертификации
```

#### Типичные проблемы:

| Проблема | Причина | Решение |
|----------|---------|---------|
| JWT validation fails | Wrong public key | Fetch from JWKS dynamically |
| Nonce replay | Same nonce used | Store nonces, check uniqueness |
| State mismatch | Lost in redirects | Secure cookie or session |
| Timezone issues | Different servers | Always UTC |
| Redirect loop | Misconfigured URLs | Validate all URLs |
| CORS errors | iframe restrictions | Proper headers |

#### Инструменты для тестирования:

```
IMS Certification Suite:
• LTI Reference Implementation (RI)
• Conformance Test Suite
• Certification badge

Testing Tools:
• https://lti-ri.imsglobal.org/ (IMS Reference Implementation)
• SALT (Sakai LTI testing tool)
• Canvas LTI test tool

Open Source LTI Tools for testing:
• H5P (interactive content)
• Tsugi (LTI tool framework)
```

---

### 6.7 Источники для дальнейшего изучения

#### Официальные спецификации:

| Ресурс | Ссылка | Описание |
|--------|--------|----------|
| **LTI 1.3 Core** | [imsglobal.org/spec/lti](https://www.imsglobal.org/spec/lti/v1p3/) | Основная спецификация |
| **LTI Advantage** | [imsglobal.org/lti-advantage](https://www.imsglobal.org/lti-advantage) | AGS, NRPS, Deep Linking |
| **Security Framework** | [imsglobal.org/spec/security](https://www.imsglobal.org/spec/security/v1p0/) | OAuth 2.0, JWT |
| **IMS Certification** | [imscert.org](https://site.imsglobal.org/certifications) | Сертификация |

#### Библиотеки:

| Язык | Библиотека | Ссылка |
|------|------------|--------|
| **Go** | go-lti | [github.com/lestrrat-go/jwx](https://github.com/lestrrat-go/jwx) (JWT) |
| **Python** | pylti1.3 | [github.com/dmitry-viskov/pylti1.3](https://github.com/dmitry-viskov/pylti1.3) |
| **PHP** | lti-1-3-php-library | [github.com/IMSGlobal](https://github.com/IMSGlobal/lti-1-3-php-library) |
| **Node.js** | ltijs | [github.com/Cvmcosta/ltijs](https://github.com/Cvmcosta/ltijs) |
| **Ruby** | lti-1.3 | [github.com/instructure/lti-1.3](https://github.com/instructure/lti-1.3) |

#### Примеры реализаций:

| Проект | Описание | Ссылка |
|--------|----------|--------|
| **Canvas LMS** | Full LTI 1.3 implementation | [github.com/instructure/canvas-lms](https://github.com/instructure/canvas-lms) |
| **Moodle** | LTI provider + consumer | [github.com/moodle/moodle](https://github.com/moodle/moodle) |
| **Tsugi** | LTI tool framework | [github.com/tsugiproject/tsugi](https://github.com/tsugiproject/tsugi) |
| **IMS RI** | Reference Implementation | [lti-ri.imsglobal.org](https://lti-ri.imsglobal.org/) |

#### Образовательные ресурсы:

```
📚 Обучающие материалы:

• IMS Global Learning Impact Leadership Institute
• "LTI 1.3 and LTI Advantage" - IMS webinar series
• Claude Ostyn's LTI blog posts
• Blackboard/Anthology LTI documentation
• Canvas LTI developer documentation

🎥 Видео:
• "Understanding LTI 1.3" - IMS Global YouTube
• "Implementing LTI Advantage" - conference talks
• "LTI Security Best Practices" - EDUCAUSE
```

---

### 6.8 Чек-лист реализации

```
Phase 1: Infrastructure (Day 1-5)
□ Database schema (tools, placements, launches, keys)
□ RSA key pair generation and storage
□ JWKS endpoint (/lti/jwks)
□ JWT library setup
□ Configuration management

Phase 2: OIDC Flow (Day 6-12)
□ Login initiation endpoint
□ Authorization endpoint
□ id_token generation
□ State/nonce management
□ Launch session tracking
□ Error handling

Phase 3: Tool Management (Day 13-18)
□ Admin UI: Add/edit/delete tools
□ Tool configuration validation
□ Per-tenant tool settings
□ Course-level tool placements
□ Custom parameters

Phase 4: LTI Advantage - AGS (Day 19-24)
□ Line items CRUD
□ Scores endpoint
□ Results endpoint
□ Grade sync to gradebook
□ OAuth 2.0 client credentials

Phase 5: LTI Advantage - NRPS & DL (Day 25-30)
□ Memberships endpoint
□ Privacy filtering
□ Deep Linking launch
□ Deep Linking response handling
□ Content item storage

Phase 6: Frontend (Day 31-36)
□ Launch button component
□ Add tool dialog (admin)
□ Tool list management
□ Deep Linking picker
□ Course tool settings

Phase 7: Testing & Certification (Day 37-45)
□ Unit tests
□ Integration tests
□ IMS Reference Implementation testing
□ Conformance suite
□ Bug fixes
□ Documentation
□ IMS Certification submission
```

---

## 7. xAPI и Learning Record Store (LRS)

### 7.1 Определение

**xAPI (Experience API, также известный как Tin Can API)** — это современный стандарт для отслеживания и записи любого учебного опыта. В отличие от SCORM, который ограничен веб-контентом внутри LMS, xAPI может фиксировать обучение везде: мобильные приложения, симуляторы, VR/AR, реальные события, социальное обучение и многое другое.

**Learning Record Store (LRS)** — это база данных, специально разработанная для хранения xAPI statements (записей об обучении). LRS может быть standalone или встроен в LMS.

**Ключевая идея:** "Track anything, anywhere" — отслеживайте любой учебный опыт в любом контексте, не только внутри LMS.

#### Сравнение SCORM vs xAPI:

| Аспект | SCORM | xAPI |
|--------|-------|------|
| **Контекст** | Только в LMS | Любое место и устройство |
| **Данные** | Ограниченный набор | Любые данные |
| **Offline** | Нет | Да, с синхронизацией |
| **Mobile** | Ограничено | Полная поддержка |
| **Детализация** | Урок/тест | Любое микро-действие |
| **Analytics** | Базовые | Продвинутые |
| **Год создания** | 2001 | 2013 |

#### Архитектура xAPI:

```
┌─────────────────────────────────────────────────────────────────┐
│                      xAPI Architecture                           │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                 Learning Activities                      │    │
│  │                                                          │    │
│  │  ┌─────────┐ ┌─────────┐ ┌─────────┐ ┌─────────┐       │    │
│  │  │   LMS   │ │ Mobile  │ │   VR    │ │Simulator│       │    │
│  │  │ Courses │ │  App    │ │Training │ │         │       │    │
│  │  └────┬────┘ └────┬────┘ └────┬────┘ └────┬────┘       │    │
│  │       │           │           │           │             │    │
│  │  ┌────┴───┐ ┌─────┴────┐ ┌────┴────┐ ┌────┴────┐       │    │
│  │  │ E-book │ │Classroom │ │  Game   │ │  Job    │       │    │
│  │  │Reading │ │  ILT     │ │ Based   │ │  Task   │       │    │
│  │  └────┬───┘ └────┬─────┘ └────┬────┘ └────┬────┘       │    │
│  │       │          │            │           │             │    │
│  └───────┼──────────┼────────────┼───────────┼─────────────┘    │
│          │          │            │           │                   │
│          └──────────┴─────┬──────┴───────────┘                   │
│                           │                                       │
│                           ▼                                       │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │              xAPI Statements (JSON)                      │    │
│  │                                                          │    │
│  │  {                                                       │    │
│  │    "actor": { "name": "Иван", "mbox": "ivan@mail.ru" },│    │
│  │    "verb": { "id": "completed" },                       │    │
│  │    "object": { "id": "course/123", "name": "Python" }   │    │
│  │  }                                                       │    │
│  └─────────────────────────────────────────────────────────┘    │
│                           │                                       │
│                           ▼                                       │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │              Learning Record Store (LRS)                 │    │
│  │                                                          │    │
│  │  • Stores all xAPI statements                           │    │
│  │  • Query API for analytics                              │    │
│  │  • Statement forwarding                                 │    │
│  │  • Aggregation & reporting                              │    │
│  │                                                          │    │
│  │  Options:                                                │    │
│  │  ├─ Embedded LRS (в нашей LMS)                          │    │
│  │  ├─ Learning Locker (open source)                       │    │
│  │  ├─ Watershed (enterprise)                              │    │
│  │  └─ SCORM Cloud LRS                                     │    │
│  └─────────────────────────────────────────────────────────┘    │
│                           │                                       │
│                           ▼                                       │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                Analytics & Reporting                     │    │
│  │                                                          │    │
│  │  • Learning analytics dashboards                        │    │
│  │  • Competency tracking                                  │    │
│  │  • Learning paths optimization                          │    │
│  │  • ROI measurement                                      │    │
│  │  • Predictive analytics                                 │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

#### Структура xAPI Statement:

```json
{
  "id": "12345678-1234-5678-1234-567812345678",
  "timestamp": "2024-01-15T10:30:00.000Z",
  
  "actor": {
    "objectType": "Agent",
    "name": "Иван Иванов",
    "mbox": "mailto:ivan@university.edu",
    "account": {
      "homePage": "https://lms.university.edu",
      "name": "user_12345"
    }
  },
  
  "verb": {
    "id": "http://adlnet.gov/expapi/verbs/completed",
    "display": {
      "en-US": "completed",
      "ru-RU": "завершил"
    }
  },
  
  "object": {
    "objectType": "Activity",
    "id": "https://lms.university.edu/courses/cs101/module/5",
    "definition": {
      "type": "http://adlnet.gov/expapi/activities/module",
      "name": {
        "en-US": "Module 5: Data Structures",
        "ru-RU": "Модуль 5: Структуры данных"
      },
      "description": {
        "en-US": "Introduction to arrays, lists, and trees"
      }
    }
  },
  
  "result": {
    "score": {
      "scaled": 0.85,
      "raw": 85,
      "min": 0,
      "max": 100
    },
    "success": true,
    "completion": true,
    "duration": "PT1H30M"
  },
  
  "context": {
    "registration": "course-enrollment-uuid",
    "instructor": {
      "name": "Dr. Петрова",
      "mbox": "mailto:petrova@university.edu"
    },
    "team": {
      "name": "Группа CS-101-A",
      "objectType": "Group"
    },
    "contextActivities": {
      "parent": [{
        "id": "https://lms.university.edu/courses/cs101",
        "definition": {
          "name": { "en-US": "CS101: Introduction to Computer Science" }
        }
      }],
      "grouping": [{
        "id": "https://lms.university.edu/programs/bachelor-cs"
      }]
    },
    "platform": "PhD Student Portal",
    "language": "ru-RU",
    "extensions": {
      "https://lms.university.edu/ext/device": "mobile",
      "https://lms.university.edu/ext/location": "library"
    }
  },
  
  "authority": {
    "objectType": "Agent",
    "name": "PhD Student Portal LRS",
    "mbox": "mailto:lrs@university.edu"
  }
}
```

#### Стандартные Verbs (глаголы):

| Verb ID | Значение | Использование |
|---------|----------|---------------|
| `attempted` | Попытался | Начал тест/задание |
| `completed` | Завершил | Закончил модуль/курс |
| `passed` | Сдал | Успешно прошел тест |
| `failed` | Не сдал | Не прошел тест |
| `answered` | Ответил | Ответил на вопрос |
| `experienced` | Ознакомился | Просмотрел контент |
| `interacted` | Взаимодействовал | Клик, scroll, hover |
| `launched` | Запустил | Открыл приложение/модуль |
| `progressed` | Продвинулся | Прогресс в курсе |
| `scored` | Получил оценку | Результат теста |
| `mastered` | Освоил | Достиг компетенции |
| `commented` | Прокомментировал | Добавил комментарий |
| `shared` | Поделился | Sharing контента |
| `asked` | Спросил | Задал вопрос |
| `attended` | Присутствовал | Посетил занятие |

#### Activity Types (типы активностей):

```
Стандартные типы активностей:
├── course          - Курс
├── module          - Модуль
├── lesson          - Урок
├── assessment      - Оценивание
├── question        - Вопрос
├── interaction     - Взаимодействие
├── media           - Медиа (видео, аудио)
├── simulation      - Симуляция
├── meeting         - Встреча
├── performance     - Практическое задание
├── file            - Файл
└── link            - Ссылка

Кастомные типы (наши):
├── phd/dissertation-defense    - Защита диссертации
├── phd/publication-submitted   - Статья подана
├── phd/advisor-meeting         - Встреча с научруком
├── phd/milestone-completed     - Веха пройдена
└── phd/competency-achieved     - Компетенция достигнута
```

---

### 7.2 Почему важно реализовать

#### Бизнес-причины:

| Причина | Описание | Влияние |
|---------|----------|---------|
| **Будущее e-learning** | xAPI заменяет SCORM | Long-term investment |
| **Learning Analytics** | Глубокий анализ обучения | Data-driven decisions |
| **Любой контекст** | Mobile, VR, ILT, workplace | Полная картина обучения |
| **Compliance** | Детальный audit trail | Регуляторные требования |
| **Персонализация** | Данные для adaptive learning | Лучший UX |

#### Технические причины:

```
📊 Преимущества xAPI:

1. Granularity (Детализация)
   SCORM: "Пользователь завершил курс"
   xAPI: "Пользователь смотрел видео 5:32, 
          поставил на паузу на 2:15, 
          перемотал на 3:00,
          ответил на вопрос за 45 сек,
          ошибся, посмотрел подсказку,
          ответил правильно"

2. Контекст
   SCORM: Только браузер, только LMS
   xAPI: Мобильное приложение в метро,
          VR симулятор в лаборатории,
          Classroom с QR check-in,
          Книга с NFC меткой

3. Offline Support
   SCORM: Требует постоянное соединение
   xAPI: Записывает локально, синхронизирует позже

4. Interoperability
   • Один LRS собирает данные из множества источников
   • Портфолио обучения переносится между организациями
```

#### Сравнение данных: SCORM vs xAPI

```
Что мы знаем с SCORM:
├── Студент начал курс: 10:00
├── Студент завершил курс: 11:30
├── Результат теста: 85%
└── Общее время: 1.5 часа

Что мы знаем с xAPI:
├── 10:00 - Открыл курс (mobile, из дома)
├── 10:05 - Начал видео "Введение"
├── 10:08 - Поставил на паузу (3:24)
├── 10:15 - Продолжил видео (другое устройство, desktop)
├── 10:22 - Завершил видео, engagement score: 78%
├── 10:25 - Открыл интерактивное упражнение
├── 10:27 - Попытка 1: неверно (вариант B)
├── 10:28 - Открыл подсказку
├── 10:29 - Попытка 2: верно
├── 10:30 - Перешел к следующему модулю
├── 10:45 - Приостановил (вышел из приложения)
├── 11:00 - Вернулся (в кафе, mobile)
├── 11:15 - Прошел тест: 85% за 15 мин
├── 11:20 - Просмотрел результаты, кликнул на объяснение ошибки
├── 11:25 - Поделился сертификатом в LinkedIn
└── Analytics: Engagement высокий, struggle на вопросе 3,
              learning style: visual, лучшее время: утро
```

#### ROI xAPI:

```
💰 Измеримые выгоды:

1. Улучшение курсов
   • Видим, где студенты "застревают"
   • Оптимизируем контент по данным
   • A/B тестирование материалов

2. Персонализация
   • Adaptive learning paths
   • Рекомендации на основе поведения
   • Индивидуальный темп

3. Compliance & Audit
   • Полный trace всех действий
   • Доказательство прохождения обучения
   • Готовность к аудитам

4. Предиктивная аналитика
   • Early warning для at-risk студентов
   • Prediction успешности
   • Оптимизация ресурсов
```

---

### 7.3 Что дает конечному пользователю

#### Для студентов:

| Функция | Описание | Польза |
|---------|----------|--------|
| **Полное портфолио** | Все обучение в одном месте | Резюме, карьера |
| **Offline обучение** | Учиться без интернета | Гибкость |
| **Cross-device** | Продолжить на другом устройстве | Удобство |
| **Персональные рекомендации** | Курсы на основе истории | Релевантность |
| **Бейджи и достижения** | Визуальное признание | Мотивация |
| **Progress insights** | Аналитика своего обучения | Самоанализ |

#### Для преподавателей:

| Функция | Описание | Польза |
|---------|----------|--------|
| **Детальная аналитика** | Видео engagement, время на задание | Улучшение контента |
| **Early warning** | Студенты в зоне риска | Своевременная помощь |
| **Сравнительный анализ** | Когорты, A/B тесты | Data-driven design |
| **Classroom tracking** | QR check-in, участие | Автоматизация |
| **Competency mapping** | Прогресс по компетенциям | Curriculam design |

#### Для организации:

| Функция | Описание | Польза |
|---------|----------|--------|
| **Learning analytics** | Dashboard по всей организации | Strategic decisions |
| **Compliance reporting** | Автоматические отчеты | Audit readiness |
| **ROI measurement** | Связь обучения с performance | Budget justification |
| **Skills gap analysis** | Какие навыки нужно развивать | Talent management |
| **Benchmarking** | Сравнение с индустрией | Competitive analysis |

#### Конкретные сценарии:

```
Сценарий 1: PhD студент - полная картина
├─ LMS: Прошел курс "Research Methods" (xAPI)
├─ Mobile: Читал статьи в PubMed (xAPI через extension)
├─ Zoom: Посетил семинар (xAPI webhook)
├─ Library: Заказал книгу (xAPI integration)
├─ Lab: Провел эксперимент (xAPI от лаб системы)
├─ Conference: Представил постер (manual xAPI)
├─ Publication: Статья принята (xAPI)
└─ LRS Dashboard: Полная картина прогресса PhD ✓

Сценарий 2: Корпоративное обучение
├─ E-learning: Прошел compliance курс
├─ Simulator: Практика на симуляторе
├─ On-the-job: Выполнил задание под наблюдением
├─ Mentor: Получил обратную связь
├─ Assessment: Сдал сертификацию
└─ HR Dashboard: Видит все + competency gap ✓

Сценарий 3: Blended learning
├─ Pre-class: Смотрел видео дома (xAPI)
├─ In-class: QR check-in (xAPI)
├─ In-class: Ответил на poll (xAPI)
├─ In-class: Групповая работа (facilitator xAPI)
├─ Post-class: Домашнее задание (xAPI)
└─ Instructor: Видит engagement каждого студента ✓
```

---

### 7.4 Какие пользователи предпочитают эту функцию

#### Целевые сегменты:

| Сегмент | Важность xAPI | Основные use cases |
|---------|--------------|-------------------|
| **Enterprise L&D** | Критическая | Compliance, skills tracking |
| **Healthcare** | Очень высокая | CME, competency, audit |
| **Aviation/Military** | Критическая | Simulation, certification |
| **Universities** | Высокая | Research, analytics |
| **K-12** | Средняя | Emerging adoption |

#### Отрасли с высоким adoption xAPI:

```
✈️ Авиация и оборона:
• Pilot training (simulators → xAPI)
• Maintenance certification
• FAA compliance reporting
• Military training (ADL initiative!)

🏥 Здравоохранение:
• CME (Continuing Medical Education)
• Competency-based training
• Procedure tracking
• Regulatory compliance

🏭 Производство:
• Safety training with VR
• Equipment certification
• On-the-job tracking
• Quality management

💼 Financial Services:
• Compliance training (must have audit trail)
• Certification tracking
• Performance correlation

🎓 Higher Education:
• Learning analytics research
• Competency-based education
• Micro-credentials
• Lifelong learning portfolios
```

#### Типичные требования:

```
Enterprise L&D:
1. "Нам нужен полный audit trail для compliance"
2. "Хотим tracking за пределами LMS (mobile, ILT)"
3. "Нужна интеграция с HR системой"
4. "Хотим предиктивную аналитику"

Healthcare:
1. "CME credits должны автоматически записываться"
2. "Нужен tracking competencies"
3. "Аудит: кто что прошел и когда"
4. "Simulation data в единый репозиторий"

Higher Ed:
1. "Исследования learning analytics"
2. "Micro-credentials и open badges"
3. "Portable learning record для студентов"
4. "Competency-based progression"
```

---

### 7.5 Как интегрировать в наше приложение

#### Стратегии интеграции LRS:

| Подход | Плюсы | Минусы | Рекомендация |
|--------|-------|--------|--------------|
| **Embedded LRS** | Полный контроль | Сложность разработки | ✅ Для MVP |
| **Learning Locker** | Open source, функциональный | Отдельный сервис | ✅ Production |
| **Watershed** | Enterprise features | Дорого | 🔶 Enterprise tier |
| **SCORM Cloud** | Простой, надежный | Vendor lock-in | 🔶 Quick start |

#### Рекомендуемая архитектура:

```
┌─────────────────────────────────────────────────────────────────┐
│                 xAPI Integration Architecture                    │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                   Our Platform                           │    │
│  │                                                          │    │
│  │  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐   │    │
│  │  │   Courses    │  │   Videos     │  │   Quizzes    │   │    │
│  │  │   Module     │  │   Player     │  │   Engine     │   │    │
│  │  └──────┬───────┘  └──────┬───────┘  └──────┬───────┘   │    │
│  │         │                 │                 │            │    │
│  │         └─────────────────┼─────────────────┘            │    │
│  │                           │                              │    │
│  │                           ▼                              │    │
│  │  ┌─────────────────────────────────────────────────────┐│    │
│  │  │              xAPI Statement Generator               ││    │
│  │  │                                                      ││    │
│  │  │  • Verb mapping                                     ││    │
│  │  │  • Actor resolution (user → agent)                  ││    │
│  │  │  • Context enrichment                               ││    │
│  │  │  • Statement validation                             ││    │
│  │  └─────────────────────────────────────────────────────┘│    │
│  │                           │                              │    │
│  └───────────────────────────┼──────────────────────────────┘    │
│                              │                                   │
│              ┌───────────────┼───────────────┐                   │
│              │               │               │                   │
│              ▼               ▼               ▼                   │
│  ┌───────────────┐  ┌──────────────┐  ┌──────────────────┐      │
│  │ Embedded LRS  │  │Learning Locker│  │  External LRS   │      │
│  │  (our DB)     │  │ (recommended) │  │  (customer's)   │      │
│  └───────────────┘  └──────────────┘  └──────────────────┘      │
│                              │                                   │
│                              ▼                                   │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │                   Analytics Layer                        │    │
│  │                                                          │    │
│  │  • Statement aggregation                                │    │
│  │  • Learning path analysis                               │    │
│  │  • Competency tracking                                  │    │
│  │  • Predictive models                                    │    │
│  │  • Dashboards (Metabase, custom)                       │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

#### Схема базы данных (Embedded LRS):

```sql
-- xAPI Actors (кэш для быстрого поиска)
CREATE TABLE xapi_actors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Actor identification
    actor_type VARCHAR(20) NOT NULL DEFAULT 'Agent',  -- 'Agent', 'Group'
    
    -- Inverse Functional Identifiers (IFI) - только один обязателен
    mbox VARCHAR(255),                    -- mailto:user@example.com
    mbox_sha1sum VARCHAR(64),             -- SHA1 hash of mbox
    openid VARCHAR(500),                  -- OpenID URL
    account_homepage VARCHAR(500),         -- Account home page
    account_name VARCHAR(255),            -- Account name
    
    -- Cached info
    name VARCHAR(255),
    
    -- Link to our user (если есть)
    user_id UUID REFERENCES users(id),
    tenant_id UUID REFERENCES tenants(id),
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    -- IFI uniqueness
    UNIQUE(mbox),
    UNIQUE(mbox_sha1sum),
    UNIQUE(openid),
    UNIQUE(account_homepage, account_name)
);

-- xAPI Verbs (справочник)
CREATE TABLE xapi_verbs (
    id VARCHAR(500) PRIMARY KEY,          -- Full IRI, e.g., http://adlnet.gov/expapi/verbs/completed
    display JSONB NOT NULL DEFAULT '{}',  -- {"en-US": "completed", "ru-RU": "завершил"}
    created_at TIMESTAMP DEFAULT NOW()
);

-- xAPI Activities (справочник)
CREATE TABLE xapi_activities (
    id VARCHAR(500) PRIMARY KEY,          -- Activity IRI
    definition JSONB,                      -- Activity definition
    /*
    {
      "type": "http://adlnet.gov/expapi/activities/course",
      "name": {"en-US": "Course Title"},
      "description": {"en-US": "Course description"},
      "moreInfo": "https://...",
      "extensions": {}
    }
    */
    
    -- Link to our entities (если есть)
    entity_type VARCHAR(50),              -- 'course', 'module', 'assessment'
    entity_id UUID,
    tenant_id UUID REFERENCES tenants(id),
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- xAPI Statements (основная таблица)
CREATE TABLE xapi_statements (
    id UUID PRIMARY KEY,                   -- Statement ID (from xAPI or generated)
    
    -- Stored timestamp (when we received it)
    stored TIMESTAMP DEFAULT NOW(),
    
    -- Statement timestamp (when it occurred)
    timestamp TIMESTAMP NOT NULL,
    
    -- Actor
    actor_id UUID NOT NULL REFERENCES xapi_actors(id),
    
    -- Verb
    verb_id VARCHAR(500) NOT NULL REFERENCES xapi_verbs(id),
    
    -- Object (can be Activity, Agent, SubStatement, StatementRef)
    object_type VARCHAR(50) NOT NULL DEFAULT 'Activity',
    object_activity_id VARCHAR(500) REFERENCES xapi_activities(id),
    object_agent_id UUID REFERENCES xapi_actors(id),
    object_statement_ref UUID,             -- Reference to another statement
    object_sub_statement JSONB,            -- For SubStatement type
    
    -- Result (optional)
    result JSONB,
    /*
    {
      "score": {"scaled": 0.85, "raw": 85, "min": 0, "max": 100},
      "success": true,
      "completion": true,
      "response": "user response text",
      "duration": "PT1H30M",
      "extensions": {}
    }
    */
    
    -- Context (optional)
    context JSONB,
    /*
    {
      "registration": "uuid",
      "instructor": {...actor...},
      "team": {...group...},
      "contextActivities": {
        "parent": [...],
        "grouping": [...],
        "category": [...],
        "other": [...]
      },
      "revision": "1.0",
      "platform": "PhD Portal",
      "language": "ru-RU",
      "statement": {...ref to another statement...},
      "extensions": {}
    }
    */
    
    -- Authority (who submitted this statement)
    authority_id UUID REFERENCES xapi_actors(id),
    
    -- Attachments reference
    has_attachments BOOLEAN DEFAULT false,
    
    -- Voided flag
    voided BOOLEAN DEFAULT false,
    voiding_statement_id UUID,
    
    -- Full statement JSON (for API responses)
    raw_statement JSONB NOT NULL,
    
    -- Tenant (for multi-tenancy)
    tenant_id UUID REFERENCES tenants(id),
    
    created_at TIMESTAMP DEFAULT NOW()
);

-- xAPI Attachments
CREATE TABLE xapi_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    statement_id UUID NOT NULL REFERENCES xapi_statements(id) ON DELETE CASCADE,
    
    usage_type VARCHAR(500) NOT NULL,      -- IRI describing attachment type
    display JSONB,                          -- Language map
    description JSONB,                      -- Language map
    content_type VARCHAR(255) NOT NULL,    -- MIME type
    length INTEGER NOT NULL,               -- Byte length
    sha2 VARCHAR(128) NOT NULL,            -- SHA-256 hash
    file_url VARCHAR(500),                 -- URL if not inline
    
    created_at TIMESTAMP DEFAULT NOW()
);

-- xAPI State (для сохранения состояния)
CREATE TABLE xapi_state (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    activity_id VARCHAR(500) NOT NULL,
    agent_id UUID NOT NULL REFERENCES xapi_actors(id),
    state_id VARCHAR(255) NOT NULL,
    registration UUID,
    
    content JSONB,
    content_type VARCHAR(255) DEFAULT 'application/json',
    etag VARCHAR(64),
    
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(activity_id, agent_id, state_id, registration)
);

-- xAPI Activity Profiles
CREATE TABLE xapi_activity_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    activity_id VARCHAR(500) NOT NULL,
    profile_id VARCHAR(255) NOT NULL,
    
    content JSONB,
    content_type VARCHAR(255) DEFAULT 'application/json',
    etag VARCHAR(64),
    
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(activity_id, profile_id)
);

-- xAPI Agent Profiles
CREATE TABLE xapi_agent_profiles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    agent_id UUID NOT NULL REFERENCES xapi_actors(id),
    profile_id VARCHAR(255) NOT NULL,
    
    content JSONB,
    content_type VARCHAR(255) DEFAULT 'application/json',
    etag VARCHAR(64),
    
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(agent_id, profile_id)
);

-- Indexes for common queries
CREATE INDEX idx_xapi_statements_actor ON xapi_statements(actor_id);
CREATE INDEX idx_xapi_statements_verb ON xapi_statements(verb_id);
CREATE INDEX idx_xapi_statements_activity ON xapi_statements(object_activity_id);
CREATE INDEX idx_xapi_statements_timestamp ON xapi_statements(timestamp);
CREATE INDEX idx_xapi_statements_stored ON xapi_statements(stored);
CREATE INDEX idx_xapi_statements_tenant ON xapi_statements(tenant_id);
CREATE INDEX idx_xapi_statements_context_registration ON xapi_statements((context->>'registration'));

-- GIN index for JSONB queries
CREATE INDEX idx_xapi_statements_context_gin ON xapi_statements USING GIN (context);
CREATE INDEX idx_xapi_statements_result_gin ON xapi_statements USING GIN (result);
```

#### Frontend компоненты:

```typescript
// frontend/src/lib/xapi/xAPIClient.ts

interface Actor {
  objectType?: 'Agent' | 'Group';
  name?: string;
  mbox?: string;
  account?: {
    homePage: string;
    name: string;
  };
}

interface Verb {
  id: string;
  display?: Record<string, string>;
}

interface Activity {
  objectType?: 'Activity';
  id: string;
  definition?: {
    type?: string;
    name?: Record<string, string>;
    description?: Record<string, string>;
    extensions?: Record<string, any>;
  };
}

interface Result {
  score?: {
    scaled?: number;
    raw?: number;
    min?: number;
    max?: number;
  };
  success?: boolean;
  completion?: boolean;
  response?: string;
  duration?: string;
  extensions?: Record<string, any>;
}

interface Context {
  registration?: string;
  instructor?: Actor;
  team?: Actor;
  contextActivities?: {
    parent?: Activity[];
    grouping?: Activity[];
    category?: Activity[];
    other?: Activity[];
  };
  platform?: string;
  language?: string;
  extensions?: Record<string, any>;
}

interface Statement {
  id?: string;
  actor: Actor;
  verb: Verb;
  object: Activity | Actor;
  result?: Result;
  context?: Context;
  timestamp?: string;
}

class XAPIClient {
  private endpoint: string;
  private actor: Actor;
  private queue: Statement[] = [];
  private flushInterval: number = 5000;
  private flushTimer?: NodeJS.Timeout;

  constructor(config: { endpoint: string; actor: Actor }) {
    this.endpoint = config.endpoint;
    this.actor = config.actor;
    this.startFlushTimer();
  }

  // Common verbs
  static VERBS = {
    LAUNCHED: { id: 'http://adlnet.gov/expapi/verbs/launched', display: { 'en-US': 'launched' } },
    COMPLETED: { id: 'http://adlnet.gov/expapi/verbs/completed', display: { 'en-US': 'completed' } },
    PASSED: { id: 'http://adlnet.gov/expapi/verbs/passed', display: { 'en-US': 'passed' } },
    FAILED: { id: 'http://adlnet.gov/expapi/verbs/failed', display: { 'en-US': 'failed' } },
    ANSWERED: { id: 'http://adlnet.gov/expapi/verbs/answered', display: { 'en-US': 'answered' } },
    EXPERIENCED: { id: 'http://adlnet.gov/expapi/verbs/experienced', display: { 'en-US': 'experienced' } },
    PROGRESSED: { id: 'http://adlnet.gov/expapi/verbs/progressed', display: { 'en-US': 'progressed' } },
    INTERACTED: { id: 'http://adlnet.gov/expapi/verbs/interacted', display: { 'en-US': 'interacted' } },
  };

  // Send a statement
  async send(statement: Omit<Statement, 'actor'>): Promise<string[]> {
    const fullStatement: Statement = {
      ...statement,
      actor: this.actor,
      timestamp: statement.timestamp || new Date().toISOString(),
    };

    // Add to queue
    this.queue.push(fullStatement);

    // If queue is getting large, flush immediately
    if (this.queue.length >= 10) {
      return this.flush();
    }

    return [];
  }

  // Convenience methods
  launched(activity: Activity, context?: Context) {
    return this.send({
      verb: XAPIClient.VERBS.LAUNCHED,
      object: activity,
      context,
    });
  }

  completed(activity: Activity, result?: Result, context?: Context) {
    return this.send({
      verb: XAPIClient.VERBS.COMPLETED,
      object: activity,
      result: { ...result, completion: true },
      context,
    });
  }

  progressed(activity: Activity, progress: number, context?: Context) {
    return this.send({
      verb: XAPIClient.VERBS.PROGRESSED,
      object: activity,
      result: {
        extensions: {
          'https://w3id.org/xapi/cmi5/result/extensions/progress': progress,
        },
      },
      context,
    });
  }

  answered(activity: Activity, response: string, correct: boolean, score?: number, context?: Context) {
    return this.send({
      verb: XAPIClient.VERBS.ANSWERED,
      object: activity,
      result: {
        response,
        success: correct,
        score: score !== undefined ? { raw: score } : undefined,
      },
      context,
    });
  }

  videoEvent(activity: Activity, event: 'played' | 'paused' | 'seeked' | 'completed', time: number, context?: Context) {
    const verbMap = {
      played: { id: 'https://w3id.org/xapi/video/verbs/played', display: { 'en-US': 'played' } },
      paused: { id: 'https://w3id.org/xapi/video/verbs/paused', display: { 'en-US': 'paused' } },
      seeked: { id: 'https://w3id.org/xapi/video/verbs/seeked', display: { 'en-US': 'seeked' } },
      completed: { id: 'http://adlnet.gov/expapi/verbs/completed', display: { 'en-US': 'completed' } },
    };

    return this.send({
      verb: verbMap[event],
      object: activity,
      result: {
        extensions: {
          'https://w3id.org/xapi/video/extensions/time': time,
        },
      },
      context,
    });
  }

  // Flush queue to server
  async flush(): Promise<string[]> {
    if (this.queue.length === 0) return [];

    const statements = [...this.queue];
    this.queue = [];

    try {
      const response = await fetch(`${this.endpoint}/statements`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-Experience-API-Version': '1.0.3',
        },
        body: JSON.stringify(statements),
      });

      if (!response.ok) {
        // Put statements back in queue for retry
        this.queue = [...statements, ...this.queue];
        throw new Error(`Failed to send statements: ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error('xAPI flush failed:', error);
      // Keep statements for retry
      this.queue = [...statements, ...this.queue];
      throw error;
    }
  }

  private startFlushTimer() {
    this.flushTimer = setInterval(() => {
      this.flush().catch(() => {});
    }, this.flushInterval);
  }

  destroy() {
    if (this.flushTimer) {
      clearInterval(this.flushTimer);
    }
    this.flush().catch(() => {});
  }
}

export { XAPIClient, type Statement, type Actor, type Activity, type Result, type Context };
```

```typescript
// frontend/src/hooks/useXAPI.ts
import { useEffect, useRef, useCallback } from 'react';
import { XAPIClient, Activity, Result, Context } from '@/lib/xapi/xAPIClient';
import { useAuth } from '@/hooks/useAuth';

export function useXAPI() {
  const { user } = useAuth();
  const clientRef = useRef<XAPIClient | null>(null);

  useEffect(() => {
    if (user) {
      clientRef.current = new XAPIClient({
        endpoint: '/api/v1/xapi',
        actor: {
          objectType: 'Agent',
          name: `${user.firstName} ${user.lastName}`,
          account: {
            homePage: window.location.origin,
            name: user.id,
          },
        },
      });
    }

    return () => {
      clientRef.current?.destroy();
    };
  }, [user]);

  const trackLaunched = useCallback((activity: Activity, context?: Context) => {
    return clientRef.current?.launched(activity, context);
  }, []);

  const trackCompleted = useCallback((activity: Activity, result?: Result, context?: Context) => {
    return clientRef.current?.completed(activity, result, context);
  }, []);

  const trackProgress = useCallback((activity: Activity, progress: number, context?: Context) => {
    return clientRef.current?.progressed(activity, progress, context);
  }, []);

  const trackAnswer = useCallback((
    activity: Activity, 
    response: string, 
    correct: boolean, 
    score?: number, 
    context?: Context
  ) => {
    return clientRef.current?.answered(activity, response, correct, score, context);
  }, []);

  const trackVideoEvent = useCallback((
    activity: Activity, 
    event: 'played' | 'paused' | 'seeked' | 'completed', 
    time: number, 
    context?: Context
  ) => {
    return clientRef.current?.videoEvent(activity, event, time, context);
  }, []);

  return {
    trackLaunched,
    trackCompleted,
    trackProgress,
    trackAnswer,
    trackVideoEvent,
    client: clientRef.current,
  };
}
```

```typescript
// frontend/src/components/Analytics/LearningAnalyticsDashboard.tsx
import React, { useEffect, useState } from 'react';
import { 
  BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer,
  LineChart, Line, PieChart, Pie, Cell 
} from 'recharts';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs';

interface AnalyticsData {
  activityByVerb: { verb: string; count: number }[];
  activityByDay: { date: string; count: number }[];
  completionRates: { course: string; rate: number }[];
  engagementScore: number;
  totalStatements: number;
  uniqueActivities: number;
  avgSessionDuration: string;
}

export function LearningAnalyticsDashboard({ userId, courseId }: { userId?: string; courseId?: string }) {
  const [data, setData] = useState<AnalyticsData | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchAnalytics() {
      const params = new URLSearchParams();
      if (userId) params.append('userId', userId);
      if (courseId) params.append('courseId', courseId);

      const response = await fetch(`/api/v1/xapi/analytics?${params}`);
      const analyticsData = await response.json();
      setData(analyticsData);
      setLoading(false);
    }

    fetchAnalytics();
  }, [userId, courseId]);

  if (loading) {
    return <div className="flex items-center justify-center h-64">Loading...</div>;
  }

  if (!data) {
    return <div>No data available</div>;
  }

  const COLORS = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#8884D8'];

  return (
    <div className="space-y-6">
      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-500">Total Activities</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{data.totalStatements.toLocaleString()}</div>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-500">Unique Content</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{data.uniqueActivities}</div>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-500">Engagement Score</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{data.engagementScore}%</div>
          </CardContent>
        </Card>
        
        <Card>
          <CardHeader className="pb-2">
            <CardTitle className="text-sm font-medium text-gray-500">Avg. Session</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-2xl font-bold">{data.avgSessionDuration}</div>
          </CardContent>
        </Card>
      </div>

      {/* Charts */}
      <Tabs defaultValue="activity">
        <TabsList>
          <TabsTrigger value="activity">Activity Over Time</TabsTrigger>
          <TabsTrigger value="verbs">Activity Types</TabsTrigger>
          <TabsTrigger value="completion">Completion Rates</TabsTrigger>
        </TabsList>

        <TabsContent value="activity">
          <Card>
            <CardHeader>
              <CardTitle>Learning Activity Over Time</CardTitle>
            </CardHeader>
            <CardContent>
              <ResponsiveContainer width="100%" height={300}>
                <LineChart data={data.activityByDay}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="date" />
                  <YAxis />
                  <Tooltip />
                  <Line type="monotone" dataKey="count" stroke="#8884d8" />
                </LineChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="verbs">
          <Card>
            <CardHeader>
              <CardTitle>Activity by Type</CardTitle>
            </CardHeader>
            <CardContent className="flex justify-center">
              <ResponsiveContainer width="100%" height={300}>
                <PieChart>
                  <Pie
                    data={data.activityByVerb}
                    dataKey="count"
                    nameKey="verb"
                    cx="50%"
                    cy="50%"
                    outerRadius={100}
                    label
                  >
                    {data.activityByVerb.map((_, index) => (
                      <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                    ))}
                  </Pie>
                  <Tooltip />
                </PieChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="completion">
          <Card>
            <CardHeader>
              <CardTitle>Course Completion Rates</CardTitle>
            </CardHeader>
            <CardContent>
              <ResponsiveContainer width="100%" height={300}>
                <BarChart data={data.completionRates} layout="vertical">
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis type="number" domain={[0, 100]} />
                  <YAxis dataKey="course" type="category" width={150} />
                  <Tooltip />
                  <Bar dataKey="rate" fill="#8884d8" />
                </BarChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  );
}
```

---

### 7.6 Сложность освоения и реализации

#### Оценка сложности:

| Аспект | Сложность | Комментарий |
|--------|-----------|-------------|
| **Спецификация** | 🟡 Средняя | Хорошо документирована |
| **Statement API** | 🟢 Низкая | Простой REST POST |
| **Query API** | 🟡 Средняя | Много параметров |
| **State API** | 🟢 Низкая | Simple CRUD |
| **Embedded LRS** | 🔴 Высокая | Много edge cases |
| **Analytics** | 🔴 Высокая | Aggregation, performance |

#### Временные оценки:

```
xAPI Statement Sending (базовое):
├── Statement format: 1 день
├── Client library: 2 дня
├── Backend endpoint: 2 дня
├── Integration в modules: 3 дня
└── Итого: 1-1.5 недели

Embedded LRS (full):
├── Database schema: 2 дня
├── Statement API: 3 дня
├── State API: 2 дня
├── Activity/Agent Profile: 2 дня
├── Query API: 4 дня
├── Voiding statements: 1 день
├── Testing: 3 дня
└── Итого: 2.5-3 недели

Learning Locker Integration:
├── Docker setup: 1 день
├── API client: 2 дня
├── Statement forwarding: 2 дня
└── Итого: 1 неделя

Analytics Dashboard:
├── Data aggregation: 3 дня
├── API endpoints: 2 дня
├── Frontend charts: 3 дня
├── Caching/performance: 2 дня
└── Итого: 1.5 недели

Общее время (full implementation): 6-8 недель
```

#### Типичные проблемы:

| Проблема | Причина | Решение |
|----------|---------|---------|
| Statement size | Слишком много данных | Ограничить extensions |
| Query performance | Много statements | Indexes, partitioning |
| Offline sync conflicts | Duplicate statements | Statement ID, voiding |
| Actor matching | Разные идентификаторы | Canonical actor resolution |
| Time drift | Клиентское время | Server timestamp |

---

### 7.7 Источники для дальнейшего изучения

#### Официальные спецификации:

| Ресурс | Ссылка | Описание |
|--------|--------|----------|
| **xAPI Spec** | [xapi.com/spec](https://github.com/adlnet/xAPI-Spec) | Официальная спецификация |
| **ADL Initiative** | [adlnet.gov](https://adlnet.gov/projects/xapi/) | Создатели xAPI |
| **xAPI Profiles** | [xapi.vocab.pub](https://xapi.vocab.pub/) | Vocabulary profiles |
| **cmi5** | [aicc.github.io/CMI-5](https://aicc.github.io/CMI-5_Spec_Current/) | xAPI profile для LMS |

#### LRS решения:

| LRS | Тип | Особенности |
|-----|-----|-------------|
| **Learning Locker** | Open Source | Популярный, хорошая документация |
| **Watershed** | Enterprise | Analytics, ADL partnership |
| **SCORM Cloud** | SaaS | Rustici, SCORM + xAPI |
| **Veracity LRS** | Enterprise | Enterprise features |
| **Yet Analytics** | Open Source | DATASIM, SQL-based |

#### Библиотеки:

| Язык | Библиотека | Ссылка |
|------|------------|--------|
| **Go** | goxapi | Custom implementation needed |
| **JavaScript** | xAPIWrapper | [github.com/adlnet/xAPIWrapper](https://github.com/adlnet/xAPIWrapper) |
| **Python** | tincan | [github.com/RusticiSoftware/TinCanPython](https://github.com/RusticiSoftware/TinCanPython) |
| **Java** | tincan-java | [github.com/RusticiSoftware/TinCanJava](https://github.com/RusticiSoftware/TinCanJava) |
| **PHP** | TinCanPHP | [github.com/RusticiSoftware/TinCanPHP](https://github.com/RusticiSoftware/TinCanPHP) |

#### Образовательные ресурсы:

```
📚 Обучающие материалы:

• xAPI.com - Introduction to xAPI
• ADL xAPI Cohort workshops
• Torrance Learning xAPI resources
• "xAPI for Instructional Designers" - O'Reilly

🎥 Видео:
• xAPI Bootcamp - ADL Initiative
• "Understanding xAPI" - DevLearn conference
• Learning Solutions Magazine webinars

🛠 Tools:
• xAPI Statement Viewer (Chrome extension)
• ADL xAPI Lab
• xAPI Statement Generator
```

---

### 7.8 Чек-лист реализации

```
Phase 1: Foundation (Day 1-5)
□ Database schema (statements, actors, verbs, activities)
□ Statement model and validation
□ UUID generation for statements
□ Timestamp handling

Phase 2: Statement API (Day 6-12)
□ POST /statements (single)
□ POST /statements (batch)
□ GET /statements (query)
□ PUT /statements (by ID)
□ Statement validation
□ Actor resolution
□ Voiding statements

Phase 3: Additional APIs (Day 13-18)
□ State API (GET/PUT/DELETE)
□ Activity Profile API
□ Agent Profile API
□ About API

Phase 4: Frontend Integration (Day 19-25)
□ xAPI client library
□ useXAPI hook
□ Video tracking
□ Quiz tracking
□ Course progress tracking
□ Offline queue

Phase 5: Analytics (Day 26-32)
□ Aggregation queries
□ Analytics API endpoints
□ Dashboard components
□ Charts (activity, completion, engagement)
□ Export functionality

Phase 6: Advanced (Day 33-40)
□ Learning Locker integration
□ Statement forwarding
□ cmi5 profile support
□ Performance optimization
□ Documentation

Phase 7: Testing (Day 41-45)
□ Conformance testing (ADL)
□ Load testing
□ Integration tests
□ Bug fixes
```

---

## 8. WCAG 2.1 AA Accessibility (Доступность)

### 8.1 Определение

**WCAG (Web Content Accessibility Guidelines)** — это международный стандарт, разработанный W3C, определяющий как сделать веб-контент доступным для людей с ограниченными возможностями. WCAG 2.1 AA — это уровень соответствия, требуемый большинством законодательств о доступности.

**Ключевая идея:** Веб должен быть доступен для всех, включая людей с нарушениями зрения, слуха, моторики и когнитивными особенностями.

#### Уровни соответствия WCAG:

| Уровень | Описание | Требование |
|---------|----------|------------|
| **A** | Минимальный | Базовая доступность |
| **AA** | Средний | Требуется законодательством ✅ |
| **AAA** | Максимальный | Идеал, сложно достичь |

#### Четыре принципа WCAG (POUR):

```
┌─────────────────────────────────────────────────────────────────┐
│                    WCAG 2.1 POUR Principles                      │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  P - Perceivable (Воспринимаемость)                     │    │
│  │                                                          │    │
│  │  Информация должна быть представлена так, чтобы          │    │
│  │  пользователи могли её воспринять                        │    │
│  │                                                          │    │
│  │  • Текстовые альтернативы для изображений               │    │
│  │  • Субтитры для видео                                   │    │
│  │  • Контраст текста                                      │    │
│  │  • Масштабирование текста                               │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  O - Operable (Управляемость)                           │    │
│  │                                                          │    │
│  │  Компоненты интерфейса должны быть управляемы           │    │
│  │                                                          │    │
│  │  • Клавиатурная навигация                               │    │
│  │  • Достаточное время                                    │    │
│  │  • Нет контента, вызывающего припадки                   │    │
│  │  • Навигация и поиск                                    │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  U - Understandable (Понятность)                        │    │
│  │                                                          │    │
│  │  Информация и управление должны быть понятны            │    │
│  │                                                          │    │
│  │  • Читаемый текст                                       │    │
│  │  • Предсказуемое поведение                              │    │
│  │  • Помощь при вводе                                     │    │
│  │  • Предотвращение ошибок                                │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                  │
│  ┌─────────────────────────────────────────────────────────┐    │
│  │  R - Robust (Надёжность)                                │    │
│  │                                                          │    │
│  │  Контент должен работать с разными технологиями         │    │
│  │                                                          │    │
│  │  • Совместимость с assistive technologies               │    │
│  │  • Валидный HTML                                        │    │
│  │  • ARIA атрибуты                                        │    │
│  └─────────────────────────────────────────────────────────┘    │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

#### Типы ограниченных возможностей:

| Категория | Примеры | Технологии | Требования |
|-----------|---------|------------|------------|
| **Зрение** | Слепота, слабое зрение, дальтонизм | Screen readers, magnifiers | Alt text, contrast, zoom |
| **Слух** | Глухота, слабый слух | Captions, transcripts | Субтитры, визуальные alerts |
| **Моторика** | Паралич, тремор, артрит | Keyboard, switch devices | Keyboard nav, large targets |
| **Когнитивные** | Дислексия, СДВГ, аутизм | Simple layout, читалки | Ясный язык, consistent UI |
| **Временные** | Сломанная рука, яркое солнце | Разные | Гибкость интерфейса |

#### Ключевые критерии WCAG 2.1 AA:

```
Perceivable (Воспринимаемость):
├── 1.1.1 Non-text Content          - Alt text для изображений
├── 1.2.1 Audio-only/Video-only     - Альтернативы для медиа
├── 1.2.2 Captions (Prerecorded)    - Субтитры для видео
├── 1.2.3 Audio Description         - Аудиоописание
├── 1.2.5 Audio Description (Pre)   - Расширенное описание
├── 1.3.1 Info and Relationships    - Семантическая разметка
├── 1.3.2 Meaningful Sequence       - Логический порядок
├── 1.3.3 Sensory Characteristics   - Не только цвет/форма
├── 1.3.4 Orientation               - Работает в любой ориентации
├── 1.3.5 Identify Input Purpose    - autocomplete атрибуты
├── 1.4.1 Use of Color              - Цвет не единственный индикатор
├── 1.4.2 Audio Control             - Контроль над звуком
├── 1.4.3 Contrast (Minimum)        - Контраст 4.5:1 для текста
├── 1.4.4 Resize Text               - Масштабирование до 200%
├── 1.4.5 Images of Text            - Избегать текста в картинках
├── 1.4.10 Reflow                   - Адаптация к 320px без scroll
├── 1.4.11 Non-text Contrast        - Контраст для UI элементов
├── 1.4.12 Text Spacing             - Настраиваемые отступы
└── 1.4.13 Content on Hover/Focus   - Управление popups

Operable (Управляемость):
├── 2.1.1 Keyboard                  - Все доступно с клавиатуры
├── 2.1.2 No Keyboard Trap          - Можно выйти Tab'ом
├── 2.1.4 Character Key Shortcuts   - Отключаемые hotkeys
├── 2.2.1 Timing Adjustable         - Регулируемое время
├── 2.2.2 Pause, Stop, Hide         - Контроль над анимацией
├── 2.3.1 Three Flashes             - Нет мигания >3 раз/сек
├── 2.4.1 Bypass Blocks             - Skip links
├── 2.4.2 Page Titled               - Осмысленные заголовки страниц
├── 2.4.3 Focus Order               - Логичный порядок фокуса
├── 2.4.4 Link Purpose (Context)    - Понятный текст ссылок
├── 2.4.5 Multiple Ways             - Несколько способов навигации
├── 2.4.6 Headings and Labels       - Описательные заголовки
├── 2.4.7 Focus Visible             - Видимый фокус
└── 2.5.1-2.5.4 Input Modalities    - Поддержка разных способов ввода

Understandable (Понятность):
├── 3.1.1 Language of Page          - lang атрибут
├── 3.1.2 Language of Parts         - lang для частей
├── 3.2.1 On Focus                  - Нет неожиданных изменений
├── 3.2.2 On Input                  - Предсказуемый ввод
├── 3.2.3 Consistent Navigation     - Одинаковая навигация
├── 3.2.4 Consistent Identification - Одинаковые названия
├── 3.3.1 Error Identification      - Идентификация ошибок
├── 3.3.2 Labels or Instructions    - Метки для полей
├── 3.3.3 Error Suggestion          - Подсказки по исправлению
└── 3.3.4 Error Prevention          - Предотвращение ошибок

Robust (Надёжность):
├── 4.1.1 Parsing                   - Валидный HTML
├── 4.1.2 Name, Role, Value         - ARIA для компонентов
└── 4.1.3 Status Messages           - Объявление статусов
```

#### Assistive Technologies:

| Технология | Описание | Популярные продукты |
|------------|----------|---------------------|
| **Screen Readers** | Озвучивание интерфейса | NVDA, JAWS, VoiceOver, TalkBack |
| **Screen Magnifiers** | Увеличение экрана | ZoomText, Windows Magnifier |
| **Voice Control** | Управление голосом | Dragon, Voice Control |
| **Switch Devices** | Переключатели для навигации | Various hardware |
| **Eye Tracking** | Управление взглядом | Tobii, EyeGaze |
| **Braille Displays** | Тактильный вывод | Various hardware |

---

### 8.2 Почему важно реализовать

#### Юридические требования:

| Страна/Регион | Закон | Требование | Штрафы |
|---------------|-------|------------|--------|
| **США** | ADA, Section 508 | WCAG 2.0 AA | Иски, штрафы |
| **ЕС** | European Accessibility Act | WCAG 2.1 AA | До 2025 обязательно |
| **Канада** | AODA | WCAG 2.0 AA | До $100K/день |
| **Великобритания** | Equality Act | WCAG 2.1 AA | Иски |
| **Австралия** | DDA | WCAG 2.0 AA | Иски |
| **Казахстан** | Закон о соц. защите | В развитии | - |

```
⚠️ Юридические риски:

США (2023-2024):
• 4,000+ ADA lawsuits против веб-сайтов
• Средний settlement: $20,000-$100,000
• Target: $6M settlement (2008)
• Domino's Pizza: Supreme Court case

Образование особенно под надзором:
• Harvard/MIT lawsuit (captions)
• Многочисленные иски против университетов
• Department of Education требует accessibility

EU (2025):
• European Accessibility Act вступает в силу
• Все digital products должны быть accessible
• Штрафы определяются странами-членами
```

#### Бизнес-причины:

| Причина | Описание | Влияние |
|---------|----------|---------|
| **Юридическое соответствие** | Избежание исков | Risk mitigation |
| **Расширение аудитории** | 15% населения с ограничениями | +15% потенциальных пользователей |
| **Тендеры** | Требование для госзакупок | Доступ к контрактам |
| **SEO** | Semantic HTML улучшает SEO | Органический трафик |
| **Качество UX** | Хороший для всех | Лучший продукт |
| **Brand reputation** | Инклюзивность | Positive image |

#### Статистика:

```
📊 Мировая статистика:

• 1.3 миллиарда людей с ограниченными возможностями (16%)
• 285 миллионов с нарушениями зрения
• 466 миллионов с нарушениями слуха
• В образовании: 1 из 5 студентов имеет disability

💰 Экономический эффект:
• $13 триллионов покупательская способность людей с disabilities
• $1.2 триллиона - только в США
• Улучшение accessibility увеличивает конверсию на 20-30%

🎓 В образовании:
• 19% студентов университетов имеют disability
• Законы требуют равного доступа к образованию
• Без accessibility = дискриминация
```

#### Влияние на всех пользователей:

```
Accessibility помогает ВСЕМ:

Субтитры:
├── Глухие/слабослышащие ← основная цель
├── Не-носители языка ← понимание
├── Шумная среда (метро, кафе) ← удобство
└── Тихая среда (библиотека, ночь) ← необходимость

Клавиатурная навигация:
├── Пользователи с моторными нарушениями ← необходимость
├── Power users ← продуктивность
├── Сломанный touchpad ← временная ситуация
└── RSI/карпальный туннель ← профилактика

Контраст:
├── Слабовидящие ← необходимость
├── Яркое солнце на экране ← частая ситуация
├── Старение (всех!) ← 40+ лет ухудшается зрение
└── Усталые глаза ← конец рабочего дня

Ясный язык:
├── Когнитивные особенности ← необходимость
├── Не-носители языка ← понимание
├── Все в состоянии стресса ← экзамены, дедлайны
└── Mobile users (маленький экран) ← удобство
```

---

### 8.3 Что дает конечному пользователю

#### Для студентов с ограниченными возможностями:

| Функция | Для кого | Польза |
|---------|----------|--------|
| **Screen reader support** | Слепые, слабовидящие | Полный доступ к контенту |
| **Субтитры/транскрипты** | Глухие, слабослышащие | Понимание видео/аудио |
| **Клавиатурная навигация** | Моторные нарушения | Управление без мыши |
| **Высокий контраст** | Слабовидящие | Читаемость |
| **Масштабирование** | Слабовидящие | Комфортный размер |
| **Ясный язык** | Дислексия, когнитивные | Понимание материала |
| **Предсказуемый UI** | Когнитивные, тревожность | Комфорт, уверенность |

#### Для всех студентов:

| Функция | Сценарий | Польза |
|---------|----------|--------|
| **Субтитры** | В транспорте, библиотеке | Просмотр без звука |
| **Транскрипты** | Поиск в лекциях | Найти нужный момент |
| **Keyboard shortcuts** | Продвинутые пользователи | Скорость |
| **Mobile-friendly** | Телефон на ходу | Гибкость |
| **Skip links** | Частый пользователь | Быстрая навигация |
| **Темная тема** | Ночью, sensitive глаза | Комфорт |

#### Для преподавателей:

| Функция | Описание | Польза |
|---------|----------|--------|
| **Accessibility checker** | Проверка контента | Создание доступных материалов |
| **Auto-captions** | Автоматические субтитры | Экономия времени |
| **Alt text prompts** | Напоминание добавить alt | Полнота контента |
| **Templates** | Доступные шаблоны | Простота |

#### Конкретные сценарии:

```
Сценарий 1: Слепой студент на курсе
├─ Screen reader читает весь интерфейс
├─ Все изображения имеют alt text
├─ Формы имеют labels
├─ Таблицы правильно структурированы
├─ PDF материалы accessible
├─ Видео имеют аудиоописание
└─ Студент полноценно учится ✓

Сценарий 2: Глухой студент
├─ Все видео имеют субтитры
├─ Живые лекции с real-time captioning
├─ Визуальные индикаторы (не только звук)
├─ Чат вместо voice features
├─ Транскрипты для подкастов
└─ Полный доступ к контенту ✓

Сценарий 3: Студент с RSI (repetitive strain injury)
├─ Полная клавиатурная навигация
├─ Голосовое управление работает
├─ Большие кликабельные области
├─ Нет drag-and-drop без альтернативы
├─ Автосохранение (меньше действий)
└─ Комфортное использование ✓

Сценарий 4: Студент с дислексией
├─ Настраиваемый шрифт (OpenDyslexic)
├─ Настраиваемый межстрочный интервал
├─ Возможность слушать текст (TTS)
├─ Ясный, простой язык
├─ Consistent layout
├─ Не перегруженный интерфейс
└─ Эффективное обучение ✓
```

---

### 8.4 Какие пользователи предпочитают эту функцию

#### Прямые бенефициары:

| Группа | % населения | Потребности |
|--------|-------------|-------------|
| **Слепые/слабовидящие** | ~4% | Screen readers, contrast, zoom |
| **Глухие/слабослышащие** | ~6% | Captions, visual alerts |
| **Моторные нарушения** | ~3% | Keyboard, large targets |
| **Когнитивные особенности** | ~5% | Simple UI, clear language |
| **Временные ограничения** | ~10% в любой момент | Гибкость интерфейса |

#### Требования по сегментам:

```
🎓 Университеты (ОБЯЗАТЕЛЬНО):
• ADA compliance (США)
• Section 508 (федеральное финансирование)
• Disability services office требует
• Risk of lawsuits
• "У нас есть студенты с disabilities"

💼 Корпорации:
• HR accessibility policies
• Diverse workforce
• Government contracts (Section 508)
• Risk management

🏛️ Государственные организации:
• Законодательные требования
• Все digital services должны быть accessible
• Аудиты accessibility

🌍 Международные организации:
• WCAG как глобальный стандарт
• EU Accessibility Act
• UN Convention on Rights of Persons with Disabilities
```

#### Типичные вопросы при продаже:

```
Вопросы от клиентов:

1. "Соответствует ли платформа WCAG 2.1 AA?"
2. "Есть ли VPAT (Voluntary Product Accessibility Template)?"
3. "Можете ли предоставить accessibility statement?"
4. "Работает ли со screen readers (NVDA, JAWS, VoiceOver)?"
5. "Есть ли субтитры для видео контента?"
6. "Поддерживается ли клавиатурная навигация?"
7. "Какой контраст у интерфейса?"
8. "Как часто проводится accessibility аудит?"

Без WCAG compliance:
• Университеты США: automatic rejection
• EU публичный сектор: невозможно
• Крупные корпорации: не пройдете procurement
```

---

### 8.5 Как интегрировать в наше приложение

#### Стратегия внедрения:

```
┌─────────────────────────────────────────────────────────────────┐
│                 Accessibility Implementation Strategy            │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Phase 1: Foundation                                             │
│  ─────────────────────                                           │
│  • Semantic HTML                                                 │
│  • ARIA landmarks                                                │
│  • Keyboard navigation                                           │
│  • Focus management                                              │
│                                                                  │
│  Phase 2: Visual Accessibility                                   │
│  ────────────────────────────                                    │
│  • Color contrast                                                │
│  • Focus indicators                                              │
│  • Responsive/reflow                                             │
│  • Text alternatives                                             │
│                                                                  │
│  Phase 3: Content Accessibility                                  │
│  ─────────────────────────────                                   │
│  • Alt text for images                                           │
│  • Captions for video                                            │
│  • Transcripts for audio                                         │
│  • Accessible documents (PDF)                                    │
│                                                                  │
│  Phase 4: Testing & Compliance                                   │
│  ────────────────────────────                                    │
│  • Automated testing                                             │
│  • Manual testing                                                │
│  • Screen reader testing                                         │
│  • User testing with disabilities                                │
│  • VPAT documentation                                            │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

#### Компонентная архитектура:

```typescript
// frontend/src/components/ui/accessible/Button.tsx
import React, { forwardRef } from 'react';
import { cva, type VariantProps } from 'class-variance-authority';
import { cn } from '@/lib/utils';

const buttonVariants = cva(
  // Base styles with accessibility
  [
    'inline-flex items-center justify-center rounded-md',
    'text-sm font-medium',
    'transition-colors',
    // Focus visible for keyboard users
    'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2',
    // Disabled state
    'disabled:pointer-events-none disabled:opacity-50',
    // Minimum touch target size (44x44px per WCAG 2.5.5)
    'min-h-[44px] min-w-[44px]',
  ],
  {
    variants: {
      variant: {
        default: [
          'bg-primary text-primary-foreground',
          'hover:bg-primary/90',
          'focus-visible:ring-primary',
        ],
        destructive: [
          'bg-destructive text-destructive-foreground',
          'hover:bg-destructive/90',
          'focus-visible:ring-destructive',
        ],
        outline: [
          'border border-input bg-background',
          'hover:bg-accent hover:text-accent-foreground',
          'focus-visible:ring-ring',
        ],
        // High contrast variant for accessibility
        highContrast: [
          'bg-black text-white',
          'hover:bg-gray-900',
          'focus-visible:ring-black focus-visible:ring-offset-2',
          'border-2 border-black',
        ],
      },
      size: {
        default: 'h-11 px-4 py-2',
        sm: 'h-10 px-3', // Still meets 44px minimum
        lg: 'h-12 px-8',
        icon: 'h-11 w-11',
      },
    },
    defaultVariants: {
      variant: 'default',
      size: 'default',
    },
  }
);

interface ButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement>,
    VariantProps<typeof buttonVariants> {
  /** Loading state - shows spinner and disables button */
  isLoading?: boolean;
  /** Screen reader text for icon-only buttons */
  srText?: string;
}

export const Button = forwardRef<HTMLButtonElement, ButtonProps>(
  ({ className, variant, size, isLoading, srText, children, disabled, ...props }, ref) => {
    return (
      <button
        ref={ref}
        className={cn(buttonVariants({ variant, size, className }))}
        disabled={disabled || isLoading}
        aria-busy={isLoading}
        aria-disabled={disabled || isLoading}
        {...props}
      >
        {isLoading && (
          <svg
            className="mr-2 h-4 w-4 animate-spin"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            aria-hidden="true"
          >
            <circle
              className="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              strokeWidth="4"
            />
            <path
              className="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
            />
          </svg>
        )}
        {children}
        {srText && <span className="sr-only">{srText}</span>}
      </button>
    );
  }
);
Button.displayName = 'Button';
```

```typescript
// frontend/src/components/ui/accessible/FormField.tsx
import React from 'react';
import { useId } from 'react';

interface FormFieldProps {
  label: string;
  error?: string;
  hint?: string;
  required?: boolean;
  children: (props: {
    id: string;
    'aria-describedby': string | undefined;
    'aria-invalid': boolean;
    'aria-required': boolean;
  }) => React.ReactNode;
}

export function FormField({ label, error, hint, required, children }: FormFieldProps) {
  const id = useId();
  const hintId = hint ? `${id}-hint` : undefined;
  const errorId = error ? `${id}-error` : undefined;
  
  const describedBy = [hintId, errorId].filter(Boolean).join(' ') || undefined;

  return (
    <div className="space-y-2">
      <label 
        htmlFor={id} 
        className="block text-sm font-medium text-gray-700"
      >
        {label}
        {required && (
          <span className="text-red-500 ml-1" aria-hidden="true">*</span>
        )}
        {required && <span className="sr-only"> (required)</span>}
      </label>
      
      {hint && (
        <p id={hintId} className="text-sm text-gray-500">
          {hint}
        </p>
      )}
      
      {children({
        id,
        'aria-describedby': describedBy,
        'aria-invalid': !!error,
        'aria-required': !!required,
      })}
      
      {error && (
        <p 
          id={errorId} 
          className="text-sm text-red-600"
          role="alert"
          aria-live="polite"
        >
          <span className="sr-only">Error: </span>
          {error}
        </p>
      )}
    </div>
  );
}

// Usage example:
// <FormField label="Email" error={errors.email} required>
//   {(props) => <input type="email" {...props} className="..." />}
// </FormField>
```

```typescript
// frontend/src/components/ui/accessible/SkipLinks.tsx
import React from 'react';

interface SkipLink {
  href: string;
  label: string;
}

interface SkipLinksProps {
  links?: SkipLink[];
}

const defaultLinks: SkipLink[] = [
  { href: '#main-content', label: 'Skip to main content' },
  { href: '#main-navigation', label: 'Skip to navigation' },
];

export function SkipLinks({ links = defaultLinks }: SkipLinksProps) {
  return (
    <nav aria-label="Skip links" className="sr-only focus-within:not-sr-only">
      <ul className="fixed top-0 left-0 z-50 flex gap-2 p-2 bg-white shadow-lg">
        {links.map((link) => (
          <li key={link.href}>
            <a
              href={link.href}
              className="block px-4 py-2 text-sm font-medium text-white bg-blue-600 rounded hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {link.label}
            </a>
          </li>
        ))}
      </ul>
    </nav>
  );
}
```

```typescript
// frontend/src/components/ui/accessible/LiveRegion.tsx
import React, { useEffect, useState } from 'react';

interface LiveRegionProps {
  message: string;
  /** 'polite' waits for user to finish, 'assertive' interrupts */
  politeness?: 'polite' | 'assertive';
  /** Whether to clear message after announcing */
  clearAfter?: number;
}

export function LiveRegion({ 
  message, 
  politeness = 'polite',
  clearAfter = 5000 
}: LiveRegionProps) {
  const [currentMessage, setCurrentMessage] = useState(message);

  useEffect(() => {
    setCurrentMessage(message);
    
    if (clearAfter && message) {
      const timer = setTimeout(() => setCurrentMessage(''), clearAfter);
      return () => clearTimeout(timer);
    }
  }, [message, clearAfter]);

  return (
    <div
      role="status"
      aria-live={politeness}
      aria-atomic="true"
      className="sr-only"
    >
      {currentMessage}
    </div>
  );
}

// Hook for easy usage
export function useLiveAnnouncement() {
  const [message, setMessage] = useState('');
  
  const announce = (text: string) => {
    setMessage(''); // Clear first to ensure re-announcement
    setTimeout(() => setMessage(text), 100);
  };

  return { message, announce };
}
```

```typescript
// frontend/src/components/ui/accessible/Modal.tsx
import React, { useEffect, useRef } from 'react';
import { createPortal } from 'react-dom';
import { X } from 'lucide-react';
import { useFocusTrap } from '@/hooks/useFocusTrap';

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  description?: string;
  children: React.ReactNode;
}

export function Modal({ isOpen, onClose, title, description, children }: ModalProps) {
  const modalRef = useRef<HTMLDivElement>(null);
  const previousFocus = useRef<HTMLElement | null>(null);
  
  // Focus trap
  useFocusTrap(modalRef, isOpen);

  // Store and restore focus
  useEffect(() => {
    if (isOpen) {
      previousFocus.current = document.activeElement as HTMLElement;
    } else if (previousFocus.current) {
      previousFocus.current.focus();
    }
  }, [isOpen]);

  // Close on Escape
  useEffect(() => {
    const handleEscape = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && isOpen) {
        onClose();
      }
    };
    
    document.addEventListener('keydown', handleEscape);
    return () => document.removeEventListener('keydown', handleEscape);
  }, [isOpen, onClose]);

  // Prevent body scroll
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
    return () => {
      document.body.style.overflow = '';
    };
  }, [isOpen]);

  if (!isOpen) return null;

  return createPortal(
    <div
      className="fixed inset-0 z-50 flex items-center justify-center"
      role="presentation"
    >
      {/* Backdrop */}
      <div
        className="absolute inset-0 bg-black/50"
        onClick={onClose}
        aria-hidden="true"
      />
      
      {/* Modal */}
      <div
        ref={modalRef}
        role="dialog"
        aria-modal="true"
        aria-labelledby="modal-title"
        aria-describedby={description ? 'modal-description' : undefined}
        className="relative z-10 w-full max-w-lg p-6 bg-white rounded-lg shadow-xl"
      >
        {/* Close button */}
        <button
          onClick={onClose}
          className="absolute top-4 right-4 p-2 rounded-full hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-blue-500"
          aria-label="Close dialog"
        >
          <X className="w-5 h-5" aria-hidden="true" />
        </button>
        
        {/* Title */}
        <h2 
          id="modal-title" 
          className="text-lg font-semibold mb-2"
        >
          {title}
        </h2>
        
        {/* Description */}
        {description && (
          <p id="modal-description" className="text-gray-600 mb-4">
            {description}
          </p>
        )}
        
        {/* Content */}
        {children}
      </div>
    </div>,
    document.body
  );
}
```

```typescript
// frontend/src/hooks/useFocusTrap.ts
import { useEffect, RefObject } from 'react';

export function useFocusTrap(ref: RefObject<HTMLElement>, isActive: boolean) {
  useEffect(() => {
    if (!isActive || !ref.current) return;

    const element = ref.current;
    const focusableElements = element.querySelectorAll<HTMLElement>(
      'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
    );
    
    const firstElement = focusableElements[0];
    const lastElement = focusableElements[focusableElements.length - 1];

    // Focus first element
    firstElement?.focus();

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key !== 'Tab') return;

      if (e.shiftKey) {
        // Shift + Tab
        if (document.activeElement === firstElement) {
          e.preventDefault();
          lastElement?.focus();
        }
      } else {
        // Tab
        if (document.activeElement === lastElement) {
          e.preventDefault();
          firstElement?.focus();
        }
      }
    };

    element.addEventListener('keydown', handleKeyDown);
    return () => element.removeEventListener('keydown', handleKeyDown);
  }, [ref, isActive]);
}
```

```typescript
// frontend/src/components/VideoPlayer/AccessibleVideoPlayer.tsx
import React, { useRef, useState } from 'react';
import { Play, Pause, Volume2, VolumeX, Settings, Subtitles } from 'lucide-react';

interface Caption {
  src: string;
  srclang: string;
  label: string;
}

interface AccessibleVideoPlayerProps {
  src: string;
  poster?: string;
  title: string;
  captions?: Caption[];
  audioDescription?: string;
  transcript?: string;
}

export function AccessibleVideoPlayer({
  src,
  poster,
  title,
  captions = [],
  audioDescription,
  transcript,
}: AccessibleVideoPlayerProps) {
  const videoRef = useRef<HTMLVideoElement>(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [isMuted, setIsMuted] = useState(false);
  const [showCaptions, setShowCaptions] = useState(true);
  const [showTranscript, setShowTranscript] = useState(false);
  const [currentTime, setCurrentTime] = useState(0);
  const [duration, setDuration] = useState(0);

  const togglePlay = () => {
    if (videoRef.current) {
      if (isPlaying) {
        videoRef.current.pause();
      } else {
        videoRef.current.play();
      }
      setIsPlaying(!isPlaying);
    }
  };

  const toggleMute = () => {
    if (videoRef.current) {
      videoRef.current.muted = !isMuted;
      setIsMuted(!isMuted);
    }
  };

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = Math.floor(seconds % 60);
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  return (
    <div className="space-y-4">
      {/* Video container */}
      <div className="relative bg-black rounded-lg overflow-hidden">
        <video
          ref={videoRef}
          src={src}
          poster={poster}
          className="w-full"
          onTimeUpdate={(e) => setCurrentTime(e.currentTarget.currentTime)}
          onLoadedMetadata={(e) => setDuration(e.currentTarget.duration)}
          onPlay={() => setIsPlaying(true)}
          onPause={() => setIsPlaying(false)}
          aria-label={title}
        >
          {/* Captions tracks */}
          {captions.map((caption, index) => (
            <track
              key={caption.srclang}
              kind="captions"
              src={caption.src}
              srcLang={caption.srclang}
              label={caption.label}
              default={index === 0 && showCaptions}
            />
          ))}
          
          {/* Audio description track */}
          {audioDescription && (
            <track
              kind="descriptions"
              src={audioDescription}
              srcLang="en"
              label="Audio Description"
            />
          )}
        </video>

        {/* Controls */}
        <div 
          className="absolute bottom-0 left-0 right-0 bg-gradient-to-t from-black/80 to-transparent p-4"
          role="toolbar"
          aria-label="Video controls"
        >
          {/* Progress bar */}
          <div className="mb-2">
            <label htmlFor="video-progress" className="sr-only">
              Video progress
            </label>
            <input
              id="video-progress"
              type="range"
              min={0}
              max={duration}
              value={currentTime}
              onChange={(e) => {
                if (videoRef.current) {
                  videoRef.current.currentTime = Number(e.target.value);
                }
              }}
              className="w-full h-1 bg-gray-400 rounded-full appearance-none cursor-pointer"
              aria-valuetext={`${formatTime(currentTime)} of ${formatTime(duration)}`}
            />
          </div>

          <div className="flex items-center gap-2">
            {/* Play/Pause */}
            <button
              onClick={togglePlay}
              className="p-2 text-white hover:bg-white/20 rounded-full focus:outline-none focus:ring-2 focus:ring-white"
              aria-label={isPlaying ? 'Pause' : 'Play'}
            >
              {isPlaying ? (
                <Pause className="w-6 h-6" aria-hidden="true" />
              ) : (
                <Play className="w-6 h-6" aria-hidden="true" />
              )}
            </button>

            {/* Mute/Unmute */}
            <button
              onClick={toggleMute}
              className="p-2 text-white hover:bg-white/20 rounded-full focus:outline-none focus:ring-2 focus:ring-white"
              aria-label={isMuted ? 'Unmute' : 'Mute'}
            >
              {isMuted ? (
                <VolumeX className="w-6 h-6" aria-hidden="true" />
              ) : (
                <Volume2 className="w-6 h-6" aria-hidden="true" />
              )}
            </button>

            {/* Time display */}
            <span className="text-white text-sm" aria-live="off">
              {formatTime(currentTime)} / {formatTime(duration)}
            </span>

            <div className="flex-1" />

            {/* Captions toggle */}
            {captions.length > 0 && (
              <button
                onClick={() => setShowCaptions(!showCaptions)}
                className={`p-2 rounded-full focus:outline-none focus:ring-2 focus:ring-white ${
                  showCaptions ? 'bg-white/30 text-white' : 'text-white hover:bg-white/20'
                }`}
                aria-label={showCaptions ? 'Hide captions' : 'Show captions'}
                aria-pressed={showCaptions}
              >
                <Subtitles className="w-6 h-6" aria-hidden="true" />
              </button>
            )}
          </div>
        </div>
      </div>

      {/* Transcript toggle */}
      {transcript && (
        <div>
          <button
            onClick={() => setShowTranscript(!showTranscript)}
            className="text-blue-600 hover:underline focus:outline-none focus:ring-2 focus:ring-blue-500 rounded"
            aria-expanded={showTranscript}
            aria-controls="video-transcript"
          >
            {showTranscript ? 'Hide transcript' : 'Show transcript'}
          </button>
          
          {showTranscript && (
            <div
              id="video-transcript"
              className="mt-4 p-4 bg-gray-50 rounded-lg max-h-64 overflow-y-auto"
            >
              <h3 className="font-medium mb-2">Transcript</h3>
              <div className="prose prose-sm">
                {transcript}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
```

#### CSS для доступности:

```css
/* frontend/src/styles/accessibility.css */

/* Focus visible - only show focus ring for keyboard users */
:focus:not(:focus-visible) {
  outline: none;
}

:focus-visible {
  outline: 2px solid var(--focus-ring-color, #2563eb);
  outline-offset: 2px;
}

/* High contrast mode support */
@media (prefers-contrast: high) {
  :root {
    --bg-primary: #ffffff;
    --text-primary: #000000;
    --border-color: #000000;
  }

  button,
  input,
  select,
  textarea {
    border: 2px solid #000000;
  }
}

/* Reduced motion */
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
    scroll-behavior: auto !important;
  }
}

/* Screen reader only class */
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  margin: -1px;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  white-space: nowrap;
  border: 0;
}

/* But visible when focused (for skip links) */
.sr-only-focusable:focus,
.sr-only-focusable:focus-within {
  position: static;
  width: auto;
  height: auto;
  padding: inherit;
  margin: inherit;
  overflow: visible;
  clip: auto;
  white-space: normal;
}

/* Minimum touch target size */
button,
[role="button"],
input[type="checkbox"],
input[type="radio"],
a {
  min-height: 44px;
  min-width: 44px;
}

/* Text spacing - allow users to override */
body {
  line-height: 1.5;
  letter-spacing: 0.12em;
  word-spacing: 0.16em;
}

p {
  margin-bottom: 2em;
}

/* Minimum contrast for text */
body {
  color: #1f2937; /* 7.5:1 contrast on white */
}

/* Link styling - not just color */
a {
  text-decoration: underline;
}

a:hover {
  text-decoration-thickness: 2px;
}

/* Error states - not just color */
.error-field {
  border-color: #dc2626;
  border-width: 2px;
  /* Also has error icon */
}

.error-field::before {
  content: "⚠️";
  margin-right: 0.5em;
}
```

---

### 8.6 Сложность освоения и реализации

#### Оценка сложности:

| Аспект | Сложность | Комментарий |
|--------|-----------|-------------|
| **Semantic HTML** | 🟢 Низкая | Базовые знания HTML |
| **ARIA** | 🟡 Средняя | Много атрибутов, легко ошибиться |
| **Keyboard Navigation** | 🟡 Средняя | Focus management tricky |
| **Screen Reader Testing** | 🔴 Высокая | Нужен опыт |
| **Video Accessibility** | 🟡 Средняя | Captions, descriptions |
| **Document Accessibility** | 🔴 Высокая | PDF accessibility сложно |
| **Color/Contrast** | 🟢 Низкая | Инструменты помогают |

#### Временные оценки:

```
Базовая accessibility (новый проект):
├── Semantic HTML: 2-3 дня
├── ARIA landmarks: 1 день
├── Keyboard navigation: 3-4 дня
├── Focus management: 2-3 дня
├── Skip links: 1 день
├── Form accessibility: 2-3 дня
├── Color contrast fix: 2 дня
└── Итого: 2-3 недели

Remediation существующего проекта:
├── Audit: 1-2 недели
├── Fixes: 4-8 недель (зависит от размера)
├── Testing: 2 недели
└── Итого: 7-12 недель

Content accessibility:
├── Alt text process: 1 неделя (setup)
├── Video captions: Ongoing (per video)
├── Transcripts: Ongoing
├── PDF remediation: 1-2 дня/документ
└── Итого: Continuous effort

VPAT documentation:
├── Initial VPAT: 2-3 недели
├── Updates: 1 неделя/квартал
└── Итого: Ongoing
```

#### Инструменты тестирования:

| Инструмент | Тип | Использование |
|------------|-----|---------------|
| **axe DevTools** | Browser extension | Автоматическое тестирование |
| **WAVE** | Browser extension | Визуальные индикаторы |
| **Lighthouse** | Chrome built-in | Accessibility score |
| **NVDA** | Screen reader (free) | Manual testing |
| **VoiceOver** | Screen reader (Mac) | Manual testing |
| **JAWS** | Screen reader (paid) | Enterprise testing |
| **Pa11y** | CLI tool | CI/CD integration |
| **jest-axe** | Jest matcher | Unit tests |

---

### 8.7 Источники для дальнейшего изучения

#### Официальные спецификации:

| Ресурс | Ссылка | Описание |
|--------|--------|----------|
| **WCAG 2.1** | [w3.org/WAI/WCAG21](https://www.w3.org/WAI/WCAG21/quickref/) | Quick Reference |
| **ARIA** | [w3.org/WAI/ARIA](https://www.w3.org/WAI/ARIA/apg/) | Authoring Practices |
| **WAI Tutorials** | [w3.org/WAI/tutorials](https://www.w3.org/WAI/tutorials/) | Практические гайды |
| **Understanding WCAG** | [w3.org/WAI/WCAG21/Understanding](https://www.w3.org/WAI/WCAG21/Understanding/) | Детальные объяснения |

#### Образовательные ресурсы:

| Ресурс | Описание | Ссылка |
|--------|----------|--------|
| **WebAIM** | Статьи, инструменты | [webaim.org](https://webaim.org/) |
| **A11y Project** | Checklist, resources | [a11yproject.com](https://www.a11yproject.com/) |
| **Deque University** | Курсы, сертификация | [dequeuniversity.com](https://dequeuniversity.com/) |
| **MDN Accessibility** | Documentation | [developer.mozilla.org](https://developer.mozilla.org/en-US/docs/Web/Accessibility) |

#### Инструменты:

| Инструмент | Тип | Ссылка |
|------------|-----|--------|
| **axe-core** | Testing library | [github.com/dequelabs/axe-core](https://github.com/dequelabs/axe-core) |
| **Pa11y** | CLI/CI tool | [pa11y.org](https://pa11y.org/) |
| **Colour Contrast Checker** | Design tool | [colourcontrast.cc](https://colourcontrast.cc/) |
| **ANDI** | Bookmarklet | [ssa.gov/accessibility/andi](https://www.ssa.gov/accessibility/andi/help/install.html) |

#### React libraries:

| Библиотека | Описание |
|------------|----------|
| **@radix-ui** | Accessible primitives |
| **@reach/ui** | Accessible components |
| **react-aria** | Adobe's accessibility hooks |
| **@headlessui/react** | Tailwind accessible components |

---

### 8.8 Чек-лист реализации

```
Phase 1: Foundation (Day 1-7)
□ Semantic HTML audit
□ ARIA landmarks (banner, main, navigation, contentinfo)
□ Heading hierarchy (h1-h6)
□ Skip links implementation
□ Language attribute (html lang)
□ Page titles

Phase 2: Keyboard & Focus (Day 8-14)
□ Tab order audit
□ Focus visible styles
□ Focus trap for modals
□ No keyboard traps
□ Custom component keyboard support
□ Shortcut keys (with disable option)

Phase 3: Forms (Day 15-21)
□ Label associations
□ Error identification
□ Error suggestions
□ Required field indication
□ Autocomplete attributes
□ Input purpose identification

Phase 4: Visual (Day 22-28)
□ Color contrast check (4.5:1 text, 3:1 UI)
□ Focus indicators (visible)
□ Non-color indicators (icons, patterns)
□ Responsive/reflow (320px)
□ Text resize (200%)
□ Text spacing support

Phase 5: Media (Day 29-35)
□ Alt text for images
□ Decorative images (alt="")
□ Complex images (long descriptions)
□ Video captions
□ Audio transcripts
□ Audio descriptions (where needed)
□ Media player accessibility

Phase 6: Testing (Day 36-45)
□ Automated testing (axe, Pa11y)
□ Manual keyboard testing
□ Screen reader testing (NVDA)
□ Screen reader testing (VoiceOver)
□ Color blindness simulation
□ User testing (with disabilities)

Phase 7: Documentation (Day 46-50)
□ VPAT creation
□ Accessibility statement
□ Known issues documentation
□ Remediation roadmap
□ Content author guidelines
```

---

## 9. OneRoster (Стандарт обмена данными)

### 9.1 Определение

**OneRoster** — это открытый стандарт IMS Global для обмена данными о классах, курсах, учениках и преподавателях между образовательными системами. Он упрощает интеграцию SIS (Student Information Systems) с LMS и другими EdTech приложениями.

**Ключевая идея:** Стандартизированный способ синхронизации roster данных (списки классов, зачисления) между системами без custom интеграций.

#### Проблема, которую решает OneRoster:

```
БЕЗ OneRoster:                        С OneRoster:
                                      
┌─────────┐                           ┌─────────┐
│   SIS   │                           │   SIS   │
└────┬────┘                           └────┬────┘
     │ Custom API #1                       │ OneRoster API
     ▼                                     ▼
┌─────────┐                           ┌──────────────┐
│  LMS A  │                           │  OneRoster   │
└─────────┘                           │   Standard   │
     │ Custom API #2                  └──────┬───────┘
     ▼                                       │
┌─────────┐                      ┌───────────┼───────────┐
│  LMS B  │                      │           │           │
└─────────┘                      ▼           ▼           ▼
     │ Custom API #3        ┌─────────┐ ┌─────────┐ ┌─────────┐
     ▼                      │  LMS A  │ │  LMS B  │ │  LMS C  │
┌─────────┐                 └─────────┘ └─────────┘ └─────────┘
│  LMS C  │                 
└─────────┘                 
                            Одна интеграция → все системы
N систем = N интеграций     
```

#### Компоненты OneRoster:

| Компонент | Описание | Использование |
|-----------|----------|---------------|
| **CSV Import/Export** | Файловый обмен | Batch sync, initial load |
| **REST API** | Real-time API | Live sync, on-demand |
| **Rostering** | Core roster data | Users, classes, enrollments |
| **Gradebook** | Оценки и результаты | Grades, categories, items |
| **Resource** | Учебные материалы | Content links |

#### Основные сущности OneRoster:

```
┌─────────────────────────────────────────────────────────────────┐
│                    OneRoster Data Model                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  Organization (Организация)                                      │
│  ├── sourcedId: "org-001"                                       │
│  ├── name: "KazNMU"                                             │
│  ├── type: "school" | "district" | "department"                 │
│  └── parent: null | Organization                                │
│                                                                  │
│  AcademicSession (Учебный период)                               │
│  ├── sourcedId: "term-2024-fall"                                │
│  ├── title: "Fall 2024"                                         │
│  ├── type: "term" | "semester" | "year"                         │
│  ├── startDate: "2024-09-01"                                    │
│  └── endDate: "2024-12-31"                                      │
│                                                                  │
│  Course (Курс/Предмет)                                          │
│  ├── sourcedId: "course-bio-101"                                │
│  ├── title: "Biology 101"                                       │
│  ├── courseCode: "BIO101"                                       │
│  ├── org: Organization                                          │
│  └── schoolYear: AcademicSession                                │
│                                                                  │
│  Class (Класс/Группа)                                           │
│  ├── sourcedId: "class-bio-101-sec-a"                           │
│  ├── title: "Biology 101 - Section A"                           │
│  ├── course: Course                                             │
│  ├── school: Organization                                       │
│  ├── terms: [AcademicSession]                                   │
│  └── classType: "homeroom" | "scheduled"                        │
│                                                                  │
│  User (Пользователь)                                            │
│  ├── sourcedId: "user-12345"                                    │
│  ├── username: "student@kaznmu.kz"                              │
│  ├── givenName: "Алмат"                                         │
│  ├── familyName: "Иванов"                                       │
│  ├── role: "student" | "teacher" | "administrator"              │
│  ├── email: "almaty@kaznmu.kz"                                  │
│  ├── orgs: [Organization]                                       │
│  └── agents: [User] (parents/guardians)                         │
│                                                                  │
│  Enrollment (Зачисление)                                        │
│  ├── sourcedId: "enroll-001"                                    │
│  ├── user: User                                                 │
│  ├── class: Class                                               │
│  ├── role: "student" | "teacher" | "aide"                       │
│  ├── primary: true | false                                      │
│  ├── beginDate: "2024-09-01"                                    │
│  └── endDate: "2024-12-31"                                      │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

#### Gradebook Extension:

```
LineItem (Задание/Оценочный элемент)
├── sourcedId: "item-quiz-1"
├── title: "Quiz 1: Cell Biology"
├── class: Class
├── category: LineItemCategory
├── dueDate: "2024-10-15"
├── assignDate: "2024-10-01"
└── resultValueMin/Max: 0, 100

LineItemCategory (Категория оценок)
├── sourcedId: "cat-quizzes"
├── title: "Quizzes"
└── weight: 20 (процент от итоговой)

Result (Результат/Оценка)
├── sourcedId: "result-001"
├── lineItem: LineItem
├── student: User
├── score: 85
├── scoreStatus: "fully graded" | "partially graded" | "exempt"
├── scoreDate: "2024-10-16"
└── comment: "Good work!"
```

#### CSV формат (для batch import):

```csv
# users.csv
sourcedId,status,dateLastModified,enabledUser,orgSourcedIds,role,username,givenName,familyName,email
user-001,active,2024-01-15T10:00:00Z,true,org-001,student,student1@kaznmu.kz,Алмат,Иванов,student1@kaznmu.kz
user-002,active,2024-01-15T10:00:00Z,true,org-001,teacher,teacher1@kaznmu.kz,Мария,Петрова,teacher1@kaznmu.kz

# classes.csv
sourcedId,status,dateLastModified,title,courseSourcedId,classCode,classType,schoolSourcedId,termSourcedIds
class-001,active,2024-01-15T10:00:00Z,Biology 101 - Section A,course-bio,BIO101-A,scheduled,org-001,term-fall-2024

# enrollments.csv
sourcedId,status,dateLastModified,classSourcedId,userSourcedId,role,primary,beginDate,endDate
enroll-001,active,2024-01-15T10:00:00Z,class-001,user-001,student,true,2024-09-01,2024-12-31
enroll-002,active,2024-01-15T10:00:00Z,class-001,user-002,teacher,true,2024-09-01,2024-12-31
```

#### REST API структура:

```
GET /ims/oneroster/v1p1/users                    # Все пользователи
GET /ims/oneroster/v1p1/users/{id}               # Конкретный пользователь
GET /ims/oneroster/v1p1/users/{id}/classes       # Классы пользователя

GET /ims/oneroster/v1p1/classes                  # Все классы
GET /ims/oneroster/v1p1/classes/{id}             # Конкретный класс
GET /ims/oneroster/v1p1/classes/{id}/students    # Студенты класса
GET /ims/oneroster/v1p1/classes/{id}/teachers    # Преподаватели класса

GET /ims/oneroster/v1p1/enrollments              # Все зачисления
GET /ims/oneroster/v1p1/courses                  # Все курсы
GET /ims/oneroster/v1p1/orgs                     # Все организации
GET /ims/oneroster/v1p1/academicSessions         # Учебные периоды

# Gradebook
GET /ims/oneroster/v1p1/classes/{id}/lineItems   # Оценочные элементы
GET /ims/oneroster/v1p1/classes/{id}/results     # Результаты класса
PUT /ims/oneroster/v1p1/results/{id}             # Обновить оценку
```

---

### 9.2 Почему важно реализовать

#### Интеграционная нагрузка без стандартов:

```
Типичный EdTech landscape школы/университета:

┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐
│   SIS   │  │   LMS   │  │ Library │  │Assessment│
│(PowerSch│  │(Canvas) │  │ System  │  │ Platform │
│  ool)   │  │         │  │         │  │          │
└────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘
     │            │            │            │
     └────────────┴────────────┴────────────┘
                       │
              Нужна синхронизация
              данных между ВСЕМИ

Без стандарта:
• N систем = N×(N-1)/2 custom интеграций
• 6 систем = 15 отдельных интеграций!
• Каждая требует поддержки
• Изменение в одной ломает другие

С OneRoster:
• N систем = N интеграций (каждая с OneRoster)
• 6 систем = 6 интеграций
• Стандарт не меняется
• Vendor-agnostic
```

#### Бизнес-причины:

| Причина | Описание | ROI |
|---------|----------|-----|
| **Сокращение интеграций** | Один стандарт вместо множества | -70% integration work |
| **Автоматизация roster** | Нет ручного ввода | -90% manual data entry |
| **Data accuracy** | Single source of truth | -80% data errors |
| **Faster deployment** | Plug & play с SIS | -60% implementation time |
| **Vendor flexibility** | Легко сменить системы | No lock-in |

#### Compliance и требования:

```
📋 Кто требует OneRoster:

США K-12:
• Большинство school districts
• State education departments
• "Must have" для тендеров

Higher Education:
• Растущее требование
• Особенно для LMS интеграций

Международно:
• UK, Australia, Canada adoption
• IMS Global certified products

Популярные SIS с OneRoster:
• PowerSchool ✓
• Infinite Campus ✓
• Skyward ✓
• Aeries ✓
• Tyler SIS ✓
• Ellucian Banner ✓ (higher ed)
```

#### Статистика:

```
📊 OneRoster adoption:

• 90%+ школьных округов США используют OneRoster-compatible SIS
• 100+ IMS certified OneRoster products
• #1 rostering standard в K-12 EdTech
• Экономия $2-5 на студента в год на data entry
• District с 50,000 студентов = $100-250K/год экономии
```

---

### 9.3 Что дает конечному пользователю

#### Для IT администраторов:

| Функция | Без OneRoster | С OneRoster |
|---------|---------------|-------------|
| **Initial setup** | Недели custom work | Часы конфигурации |
| **User provisioning** | Manual или scripts | Automatic sync |
| **Class creation** | Manual или scripts | Automatic sync |
| **Enrollment updates** | Daily manual work | Real-time sync |
| **Error handling** | Debugging custom code | Standard error format |

#### Для преподавателей:

| Функция | Описание | Польза |
|---------|----------|--------|
| **Ready rosters** | Классы уже созданы | Не тратить время на setup |
| **Correct students** | Списки синхронизированы | Нет "phantom students" |
| **Grade passback** | Оценки идут в SIS | Один ввод оценок |
| **Updated info** | Актуальные данные | Правильные emails, имена |

#### Для студентов:

| Функция | Описание | Польза |
|---------|----------|--------|
| **Instant access** | Автоматическое зачисление | Доступ с первого дня |
| **Correct classes** | Правильное расписание | Нет путаницы |
| **Unified grades** | Оценки везде одинаковы | Один источник правды |
| **Schedule changes** | Автоматические обновления | Нет ручных запросов |

#### Конкретные сценарии:

```
Сценарий 1: Начало семестра
├─ SIS: 5000 студентов зачислены на курсы
├─ OneRoster sync: Автоматически
├─ LMS: Все курсы созданы, студенты зачислены
├─ Day 1: Студенты могут учиться
└─ IT effort: 0 часов (вместо 40+)

Сценарий 2: Студент меняет группу
├─ SIS: Admin меняет enrollment
├─ OneRoster sync: Через минуты
├─ LMS: Студент видит новый класс
├─ Старый класс: Доступ удален
└─ Manual work: 0 (вместо 4 systems)

Сценарий 3: Преподаватель ставит оценки
├─ LMS: Преподаватель вводит оценки
├─ OneRoster Gradebook: Sync to SIS
├─ SIS: Оценки в официальном журнале
├─ Parent portal: Родители видят оценки
└─ Data entry: 1x вместо 3x

Сценарий 4: Mid-semester enrollment
├─ New student: Зачислен в SIS
├─ OneRoster: Immediate sync
├─ All systems: Access granted
├─ Student: Can start learning today
└─ Wait time: Minutes instead of days
```

---

### 9.4 Какие пользователи предпочитают эту функцию

#### Прямые бенефициары:

| Сегмент | Требование | Причина |
|---------|------------|---------|
| **K-12 школы (США)** | ОБЯЗАТЕЛЬНО | Все SIS поддерживают, стандарт де-факто |
| **School districts** | ОБЯЗАТЕЛЬНО | Тысячи студентов, ручной ввод невозможен |
| **Universities** | Растёт | Ellucian, Workday интеграции |
| **EdTech vendors** | ОБЯЗАТЕЛЬНО | Для продаж в K-12 |

#### Типичные вопросы при продаже:

```
Вопросы от школьных округов:

1. "Поддерживаете ли OneRoster 1.1 или 1.2?"
2. "Есть ли IMS certification?"
3. "Можете ли импортировать CSV или нужен API?"
4. "Поддерживается ли Gradebook exchange?"
5. "Как часто происходит sync?"
6. "Работаете ли с PowerSchool/Infinite Campus/Skyward?"
7. "Какой OAuth 2.0 flow используется?"

Без OneRoster support:
• K-12 market: Практически закрыт
• Тендеры: Automatic disqualification
• Manual onboarding: Неприемлемо для districts
```

---

### 9.5 Как интегрировать в наше приложение

#### Архитектура интеграции:

```
┌─────────────────────────────────────────────────────────────────┐
│                 OneRoster Integration Architecture               │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  External SIS                     Our Platform                   │
│  ┌─────────────┐                 ┌───────────────────────┐      │
│  │             │                 │                       │      │
│  │ PowerSchool │                 │  ┌─────────────────┐  │      │
│  │ Infinite    │  OneRoster API  │  │  OneRoster      │  │      │
│  │ Campus      │◄───────────────►│  │  Service        │  │      │
│  │ etc.        │                 │  │                 │  │      │
│  │             │                 │  │  • CSV Import   │  │      │
│  └─────────────┘                 │  │  • REST Client  │  │      │
│                                  │  │  • REST Server  │  │      │
│        OR                        │  │  • Sync Engine  │  │      │
│                                  │  └────────┬────────┘  │      │
│  ┌─────────────┐                 │           │           │      │
│  │  CSV Files  │                 │           ▼           │      │
│  │             │  File Upload    │  ┌─────────────────┐  │      │
│  │ • users.csv │────────────────►│  │  Data Mapper    │  │      │
│  │ • classes   │                 │  │                 │  │      │
│  │ • enrolls   │                 │  │  OneRoster →    │  │      │
│  └─────────────┘                 │  │  Internal Model │  │      │
│                                  │  └────────┬────────┘  │      │
│                                  │           │           │      │
│                                  │           ▼           │      │
│                                  │  ┌─────────────────┐  │      │
│                                  │  │  Core Services  │  │      │
│                                  │  │                 │  │      │
│                                  │  │  • Users        │  │      │
│                                  │  │  • Courses      │  │      │
│                                  │  │  • Enrollments  │  │      │
│                                  │  │  • Grades       │  │      │
│                                  │  └─────────────────┘  │      │
│                                  │                       │      │
│                                  └───────────────────────┘      │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘
```

#### Database Schema:

```sql
-- OneRoster sync configuration
CREATE TABLE oneroster_connections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    name VARCHAR(255) NOT NULL,
    connection_type VARCHAR(50) NOT NULL, -- 'csv', 'api_client', 'api_server'
    
    -- API connection settings
    base_url VARCHAR(500),
    client_id VARCHAR(255),
    client_secret_encrypted TEXT,
    oauth_token_url VARCHAR(500),
    
    -- Sync settings
    sync_frequency VARCHAR(50) DEFAULT 'daily', -- 'realtime', 'hourly', 'daily'
    last_sync_at TIMESTAMP WITH TIME ZONE,
    last_sync_status VARCHAR(50),
    last_sync_error TEXT,
    
    -- Feature flags
    sync_users BOOLEAN DEFAULT true,
    sync_classes BOOLEAN DEFAULT true,
    sync_enrollments BOOLEAN DEFAULT true,
    sync_grades BOOLEAN DEFAULT false,
    
    -- Mapping settings
    user_role_mapping JSONB DEFAULT '{}',
    org_mapping JSONB DEFAULT '{}',
    
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Mapping table: OneRoster sourcedId → Internal ID
CREATE TABLE oneroster_id_mappings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    connection_id UUID NOT NULL REFERENCES oneroster_connections(id),
    entity_type VARCHAR(50) NOT NULL, -- 'user', 'class', 'course', 'enrollment', 'org'
    sourced_id VARCHAR(255) NOT NULL,
    internal_id UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(connection_id, entity_type, sourced_id)
);

CREATE INDEX idx_oneroster_mappings_sourced ON oneroster_id_mappings(connection_id, entity_type, sourced_id);
CREATE INDEX idx_oneroster_mappings_internal ON oneroster_id_mappings(internal_id);

-- Sync log
CREATE TABLE oneroster_sync_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    connection_id UUID NOT NULL REFERENCES oneroster_connections(id),
    sync_type VARCHAR(50) NOT NULL, -- 'full', 'delta', 'manual'
    started_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE,
    status VARCHAR(50) NOT NULL, -- 'running', 'completed', 'failed'
    
    -- Statistics
    users_created INT DEFAULT 0,
    users_updated INT DEFAULT 0,
    users_deactivated INT DEFAULT 0,
    classes_created INT DEFAULT 0,
    classes_updated INT DEFAULT 0,
    enrollments_created INT DEFAULT 0,
    enrollments_removed INT DEFAULT 0,
    
    errors JSONB DEFAULT '[]',
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Academic sessions mapping
CREATE TABLE oneroster_academic_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    connection_id UUID NOT NULL REFERENCES oneroster_connections(id),
    sourced_id VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    session_type VARCHAR(50), -- 'term', 'semester', 'year'
    start_date DATE,
    end_date DATE,
    status VARCHAR(50) DEFAULT 'active',
    internal_term_id UUID REFERENCES academic_terms(id),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(connection_id, sourced_id)
);

-- Organizations mapping
CREATE TABLE oneroster_organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    connection_id UUID NOT NULL REFERENCES oneroster_connections(id),
    sourced_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    org_type VARCHAR(50), -- 'school', 'district', 'department'
    identifier VARCHAR(255),
    parent_sourced_id VARCHAR(255),
    internal_department_id UUID REFERENCES departments(id),
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    
    UNIQUE(connection_id, sourced_id)
);
```

#### Frontend - OneRoster Configuration:

```typescript
// frontend/src/components/admin/OneRosterConfig.tsx
import React, { useState } from 'react';
import { useQuery, useMutation } from '@tanstack/react-query';
import { 
  Upload, 
  RefreshCw, 
  Settings, 
  CheckCircle, 
  XCircle,
  Clock,
  Database,
  Users,
  BookOpen,
  UserPlus
} from 'lucide-react';

interface OneRosterConnection {
  id: string;
  name: string;
  connectionType: 'csv' | 'api_client' | 'api_server';
  baseUrl?: string;
  syncFrequency: string;
  lastSyncAt?: string;
  lastSyncStatus?: string;
  isActive: boolean;
  syncUsers: boolean;
  syncClasses: boolean;
  syncEnrollments: boolean;
  syncGrades: boolean;
}

interface SyncStats {
  usersCreated: number;
  usersUpdated: number;
  usersDeactivated: number;
  classesCreated: number;
  classesUpdated: number;
  enrollmentsCreated: number;
  enrollmentsRemoved: number;
}

export function OneRosterConfig() {
  const [showAddDialog, setShowAddDialog] = useState(false);
  const [selectedConnection, setSelectedConnection] = useState<string | null>(null);

  const { data: connections, isLoading } = useQuery<OneRosterConnection[]>({
    queryKey: ['oneroster-connections'],
    queryFn: () => fetch('/api/admin/oneroster/connections').then(r => r.json()),
  });

  const syncMutation = useMutation({
    mutationFn: (connectionId: string) =>
      fetch(`/api/admin/oneroster/connections/${connectionId}/sync`, {
        method: 'POST',
      }),
    onSuccess: () => {
      // Refetch connections to update status
    },
  });

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">OneRoster Integration</h2>
          <p className="text-gray-600">
            Синхронизация данных с Student Information Systems
          </p>
        </div>
        <button
          onClick={() => setShowAddDialog(true)}
          className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
        >
          + Add Connection
        </button>
      </div>

      {/* Connection Cards */}
      <div className="grid gap-4">
        {connections?.map((conn) => (
          <ConnectionCard
            key={conn.id}
            connection={conn}
            onSync={() => syncMutation.mutate(conn.id)}
            onConfigure={() => setSelectedConnection(conn.id)}
            isSyncing={syncMutation.isPending}
          />
        ))}
      </div>

      {/* Empty State */}
      {connections?.length === 0 && (
        <div className="text-center py-12 bg-gray-50 rounded-lg">
          <Database className="w-12 h-12 mx-auto text-gray-400 mb-4" />
          <h3 className="text-lg font-medium mb-2">No Connections</h3>
          <p className="text-gray-600 mb-4">
            Connect your Student Information System to enable automatic roster sync
          </p>
          <button
            onClick={() => setShowAddDialog(true)}
            className="px-4 py-2 bg-blue-600 text-white rounded-lg"
          >
            Add Your First Connection
          </button>
        </div>
      )}

      {/* Add Connection Dialog */}
      {showAddDialog && (
        <AddConnectionDialog onClose={() => setShowAddDialog(false)} />
      )}

      {/* Configure Connection Dialog */}
      {selectedConnection && (
        <ConfigureConnectionDialog
          connectionId={selectedConnection}
          onClose={() => setSelectedConnection(null)}
        />
      )}
    </div>
  );
}

function ConnectionCard({
  connection,
  onSync,
  onConfigure,
  isSyncing,
}: {
  connection: OneRosterConnection;
  onSync: () => void;
  onConfigure: () => void;
  isSyncing: boolean;
}) {
  const statusColors = {
    completed: 'text-green-600 bg-green-100',
    failed: 'text-red-600 bg-red-100',
    running: 'text-blue-600 bg-blue-100',
  };

  return (
    <div className="bg-white border rounded-lg p-6">
      <div className="flex items-start justify-between">
        <div className="flex items-center gap-4">
          <div className="p-3 bg-blue-100 rounded-lg">
            <Database className="w-6 h-6 text-blue-600" />
          </div>
          <div>
            <h3 className="font-semibold text-lg">{connection.name}</h3>
            <p className="text-gray-600 text-sm">
              {connection.connectionType === 'csv' && 'CSV Import'}
              {connection.connectionType === 'api_client' && 'API Client (Pull)'}
              {connection.connectionType === 'api_server' && 'API Server (Push)'}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          {connection.lastSyncStatus && (
            <span
              className={`px-2 py-1 rounded-full text-xs font-medium ${
                statusColors[connection.lastSyncStatus as keyof typeof statusColors]
              }`}
            >
              {connection.lastSyncStatus === 'completed' && <CheckCircle className="w-3 h-3 inline mr-1" />}
              {connection.lastSyncStatus === 'failed' && <XCircle className="w-3 h-3 inline mr-1" />}
              {connection.lastSyncStatus}
            </span>
          )}
        </div>
      </div>

      {/* Sync Features */}
      <div className="mt-4 flex flex-wrap gap-2">
        {connection.syncUsers && (
          <span className="px-2 py-1 bg-gray-100 rounded text-xs flex items-center gap-1">
            <Users className="w-3 h-3" /> Users
          </span>
        )}
        {connection.syncClasses && (
          <span className="px-2 py-1 bg-gray-100 rounded text-xs flex items-center gap-1">
            <BookOpen className="w-3 h-3" /> Classes
          </span>
        )}
        {connection.syncEnrollments && (
          <span className="px-2 py-1 bg-gray-100 rounded text-xs flex items-center gap-1">
            <UserPlus className="w-3 h-3" /> Enrollments
          </span>
        )}
        {connection.syncGrades && (
          <span className="px-2 py-1 bg-purple-100 text-purple-700 rounded text-xs">
            Gradebook
          </span>
        )}
      </div>

      {/* Last Sync Info */}
      {connection.lastSyncAt && (
        <div className="mt-4 text-sm text-gray-500 flex items-center gap-2">
          <Clock className="w-4 h-4" />
          Last synced: {new Date(connection.lastSyncAt).toLocaleString()}
        </div>
      )}

      {/* Actions */}
      <div className="mt-4 pt-4 border-t flex justify-end gap-2">
        <button
          onClick={onConfigure}
          className="px-3 py-2 text-gray-600 hover:bg-gray-100 rounded-lg flex items-center gap-2"
        >
          <Settings className="w-4 h-4" />
          Configure
        </button>
        <button
          onClick={onSync}
          disabled={isSyncing}
          className="px-3 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 flex items-center gap-2 disabled:opacity-50"
        >
          <RefreshCw className={`w-4 h-4 ${isSyncing ? 'animate-spin' : ''}`} />
          {isSyncing ? 'Syncing...' : 'Sync Now'}
        </button>
      </div>
    </div>
  );
}
```

```typescript
// frontend/src/components/admin/CSVImportWizard.tsx
import React, { useState, useCallback } from 'react';
import { useDropzone } from 'react-dropzone';
import { Upload, FileText, Check, AlertCircle, ChevronRight } from 'lucide-react';

interface CSVFile {
  name: string;
  file: File;
  status: 'pending' | 'validated' | 'error';
  rowCount?: number;
  errors?: string[];
}

const REQUIRED_FILES = [
  { key: 'orgs', name: 'orgs.csv', description: 'Organizations' },
  { key: 'users', name: 'users.csv', description: 'Users (students, teachers)' },
  { key: 'courses', name: 'courses.csv', description: 'Courses/Subjects' },
  { key: 'classes', name: 'classes.csv', description: 'Classes/Sections' },
  { key: 'enrollments', name: 'enrollments.csv', description: 'Class enrollments' },
];

const OPTIONAL_FILES = [
  { key: 'academicSessions', name: 'academicSessions.csv', description: 'Terms/Semesters' },
  { key: 'demographics', name: 'demographics.csv', description: 'User demographics' },
];

export function CSVImportWizard({ connectionId, onComplete }: { 
  connectionId: string;
  onComplete: () => void;
}) {
  const [step, setStep] = useState(1);
  const [files, setFiles] = useState<Record<string, CSVFile>>({});
  const [validationResults, setValidationResults] = useState<any>(null);
  const [importProgress, setImportProgress] = useState(0);

  const onDrop = useCallback((acceptedFiles: File[]) => {
    const newFiles = { ...files };
    
    acceptedFiles.forEach((file) => {
      const key = file.name.replace('.csv', '');
      newFiles[key] = {
        name: file.name,
        file,
        status: 'pending',
      };
    });
    
    setFiles(newFiles);
  }, [files]);

  const { getRootProps, getInputProps, isDragActive } = useDropzone({
    onDrop,
    accept: { 'text/csv': ['.csv'] },
    multiple: true,
  });

  const validateFiles = async () => {
    // Validate each file
    const formData = new FormData();
    Object.values(files).forEach((f) => {
      formData.append(f.name, f.file);
    });

    const response = await fetch(
      `/api/admin/oneroster/connections/${connectionId}/validate`,
      { method: 'POST', body: formData }
    );
    
    const results = await response.json();
    setValidationResults(results);
    
    // Update file statuses
    const updatedFiles = { ...files };
    Object.entries(results.files).forEach(([key, result]: [string, any]) => {
      if (updatedFiles[key]) {
        updatedFiles[key].status = result.valid ? 'validated' : 'error';
        updatedFiles[key].rowCount = result.rowCount;
        updatedFiles[key].errors = result.errors;
      }
    });
    setFiles(updatedFiles);
    
    if (results.valid) {
      setStep(3);
    }
  };

  const startImport = async () => {
    setStep(4);
    
    const formData = new FormData();
    Object.values(files).forEach((f) => {
      formData.append(f.name, f.file);
    });

    const response = await fetch(
      `/api/admin/oneroster/connections/${connectionId}/import`,
      { method: 'POST', body: formData }
    );

    // Handle streaming progress...
    const reader = response.body?.getReader();
    // ... progress updates
    
    onComplete();
  };

  return (
    <div className="max-w-3xl mx-auto">
      {/* Steps indicator */}
      <div className="flex items-center justify-center mb-8">
        {[1, 2, 3, 4].map((s) => (
          <React.Fragment key={s}>
            <div
              className={`w-8 h-8 rounded-full flex items-center justify-center ${
                s <= step ? 'bg-blue-600 text-white' : 'bg-gray-200 text-gray-600'
              }`}
            >
              {s < step ? <Check className="w-4 h-4" /> : s}
            </div>
            {s < 4 && (
              <div
                className={`w-16 h-1 ${s < step ? 'bg-blue-600' : 'bg-gray-200'}`}
              />
            )}
          </React.Fragment>
        ))}
      </div>

      {/* Step 1: Upload */}
      {step === 1 && (
        <div className="space-y-6">
          <h3 className="text-xl font-semibold">Upload OneRoster CSV Files</h3>
          
          <div
            {...getRootProps()}
            className={`border-2 border-dashed rounded-lg p-8 text-center cursor-pointer transition-colors ${
              isDragActive ? 'border-blue-500 bg-blue-50' : 'border-gray-300 hover:border-blue-400'
            }`}
          >
            <input {...getInputProps()} />
            <Upload className="w-12 h-12 mx-auto text-gray-400 mb-4" />
            <p className="text-lg mb-2">
              Drag & drop OneRoster CSV files here
            </p>
            <p className="text-gray-500 text-sm">
              or click to select files
            </p>
          </div>

          {/* Required files checklist */}
          <div>
            <h4 className="font-medium mb-2">Required Files:</h4>
            <div className="space-y-2">
              {REQUIRED_FILES.map((rf) => (
                <div
                  key={rf.key}
                  className={`flex items-center gap-3 p-2 rounded ${
                    files[rf.key] ? 'bg-green-50' : 'bg-gray-50'
                  }`}
                >
                  {files[rf.key] ? (
                    <Check className="w-5 h-5 text-green-600" />
                  ) : (
                    <FileText className="w-5 h-5 text-gray-400" />
                  )}
                  <span className="font-mono text-sm">{rf.name}</span>
                  <span className="text-gray-500 text-sm">- {rf.description}</span>
                </div>
              ))}
            </div>
          </div>

          <button
            onClick={() => setStep(2)}
            disabled={REQUIRED_FILES.some((rf) => !files[rf.key])}
            className="w-full py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
          >
            Continue to Validation
            <ChevronRight className="w-4 h-4" />
          </button>
        </div>
      )}

      {/* Step 2: Validate */}
      {step === 2 && (
        <div className="space-y-6">
          <h3 className="text-xl font-semibold">Validating Files</h3>
          
          <div className="space-y-3">
            {Object.values(files).map((f) => (
              <div
                key={f.name}
                className="flex items-center justify-between p-3 border rounded-lg"
              >
                <div className="flex items-center gap-3">
                  <FileText className="w-5 h-5 text-gray-400" />
                  <span>{f.name}</span>
                </div>
                <div>
                  {f.status === 'pending' && (
                    <span className="text-gray-500">Pending</span>
                  )}
                  {f.status === 'validated' && (
                    <span className="text-green-600 flex items-center gap-1">
                      <Check className="w-4 h-4" />
                      {f.rowCount} rows
                    </span>
                  )}
                  {f.status === 'error' && (
                    <span className="text-red-600 flex items-center gap-1">
                      <AlertCircle className="w-4 h-4" />
                      {f.errors?.length} errors
                    </span>
                  )}
                </div>
              </div>
            ))}
          </div>

          <button
            onClick={validateFiles}
            className="w-full py-3 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
          >
            Validate Files
          </button>
        </div>
      )}

      {/* Step 3: Preview */}
      {step === 3 && validationResults && (
        <div className="space-y-6">
          <h3 className="text-xl font-semibold">Import Preview</h3>
          
          <div className="bg-gray-50 rounded-lg p-6">
            <h4 className="font-medium mb-4">This import will create:</h4>
            <div className="grid grid-cols-2 gap-4">
              <div className="bg-white p-4 rounded-lg">
                <div className="text-3xl font-bold text-blue-600">
                  {validationResults.preview.users}
                </div>
                <div className="text-gray-600">Users</div>
              </div>
              <div className="bg-white p-4 rounded-lg">
                <div className="text-3xl font-bold text-blue-600">
                  {validationResults.preview.classes}
                </div>
                <div className="text-gray-600">Classes</div>
              </div>
              <div className="bg-white p-4 rounded-lg">
                <div className="text-3xl font-bold text-blue-600">
                  {validationResults.preview.enrollments}
                </div>
                <div className="text-gray-600">Enrollments</div>
              </div>
              <div className="bg-white p-4 rounded-lg">
                <div className="text-3xl font-bold text-blue-600">
                  {validationResults.preview.courses}
                </div>
                <div className="text-gray-600">Courses</div>
              </div>
            </div>
          </div>

          <button
            onClick={startImport}
            className="w-full py-3 bg-green-600 text-white rounded-lg hover:bg-green-700"
          >
            Start Import
          </button>
        </div>
      )}

      {/* Step 4: Importing */}
      {step === 4 && (
        <div className="space-y-6 text-center">
          <h3 className="text-xl font-semibold">Importing Data</h3>
          
          <div className="w-full bg-gray-200 rounded-full h-4">
            <div
              className="bg-blue-600 h-4 rounded-full transition-all"
              style={{ width: `${importProgress}%` }}
            />
          </div>
          
          <p className="text-gray-600">
            {importProgress}% complete
          </p>
        </div>
      )}
    </div>
  );
}
```

```typescript
// frontend/src/components/admin/SyncHistory.tsx
import React from 'react';
import { useQuery } from '@tanstack/react-query';
import { 
  CheckCircle, 
  XCircle, 
  Clock, 
  Users, 
  BookOpen, 
  UserPlus,
  AlertTriangle
} from 'lucide-react';

interface SyncLog {
  id: string;
  syncType: 'full' | 'delta' | 'manual';
  startedAt: string;
  completedAt?: string;
  status: 'running' | 'completed' | 'failed';
  usersCreated: number;
  usersUpdated: number;
  usersDeactivated: number;
  classesCreated: number;
  classesUpdated: number;
  enrollmentsCreated: number;
  enrollmentsRemoved: number;
  errors: { message: string; entity: string; sourcedId: string }[];
}

export function SyncHistory({ connectionId }: { connectionId: string }) {
  const { data: logs } = useQuery<SyncLog[]>({
    queryKey: ['oneroster-sync-logs', connectionId],
    queryFn: () =>
      fetch(`/api/admin/oneroster/connections/${connectionId}/logs`)
        .then(r => r.json()),
  });

  return (
    <div className="space-y-4">
      <h3 className="font-semibold text-lg">Sync History</h3>

      <div className="space-y-3">
        {logs?.map((log) => (
          <div key={log.id} className="border rounded-lg p-4">
            <div className="flex items-center justify-between mb-3">
              <div className="flex items-center gap-2">
                {log.status === 'completed' && (
                  <CheckCircle className="w-5 h-5 text-green-600" />
                )}
                {log.status === 'failed' && (
                  <XCircle className="w-5 h-5 text-red-600" />
                )}
                {log.status === 'running' && (
                  <Clock className="w-5 h-5 text-blue-600 animate-pulse" />
                )}
                <span className="font-medium capitalize">{log.syncType} Sync</span>
              </div>
              <span className="text-sm text-gray-500">
                {new Date(log.startedAt).toLocaleString()}
              </span>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-3 gap-4 text-sm">
              <div className="flex items-center gap-2">
                <Users className="w-4 h-4 text-gray-400" />
                <span>
                  +{log.usersCreated} / ~{log.usersUpdated} / -{log.usersDeactivated} users
                </span>
              </div>
              <div className="flex items-center gap-2">
                <BookOpen className="w-4 h-4 text-gray-400" />
                <span>
                  +{log.classesCreated} / ~{log.classesUpdated} classes
                </span>
              </div>
              <div className="flex items-center gap-2">
                <UserPlus className="w-4 h-4 text-gray-400" />
                <span>
                  +{log.enrollmentsCreated} / -{log.enrollmentsRemoved} enrollments
                </span>
              </div>
            </div>

            {/* Errors */}
            {log.errors.length > 0 && (
              <div className="mt-3 pt-3 border-t">
                <div className="flex items-center gap-2 text-amber-600 mb-2">
                  <AlertTriangle className="w-4 h-4" />
                  <span className="text-sm font-medium">
                    {log.errors.length} errors
                  </span>
                </div>
                <div className="text-sm text-gray-600 space-y-1">
                  {log.errors.slice(0, 3).map((err, i) => (
                    <div key={i}>
                      {err.entity} ({err.sourcedId}): {err.message}
                    </div>
                  ))}
                  {log.errors.length > 3 && (
                    <div className="text-gray-400">
                      +{log.errors.length - 3} more errors
                    </div>
                  )}
                </div>
              </div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
}
```

---

### 9.6 Сложность освоения и реализации

#### Оценка сложности:

| Аспект | Сложность | Комментарий |
|--------|-----------|-------------|
| **CSV Import** | 🟢 Низкая | Парсинг стандартных файлов |
| **REST API Client** | 🟡 Средняя | OAuth 2.0, pagination |
| **REST API Server** | 🟡 Средняя | Expose наши данные |
| **Data Mapping** | 🟡 Средняя | Маппинг на internal model |
| **Sync Logic** | 🔴 Высокая | Delta sync, conflict resolution |
| **Gradebook** | 🟡 Средняя | Bi-directional sync |

#### Временные оценки:

```
MVP (CSV Import only):
├── CSV parsing: 2-3 дня
├── Data mapping: 3-4 дня
├── Import wizard UI: 2-3 дня
├── Validation: 2 дня
├── Testing: 3 дня
└── Итого: 2-3 недели

Full Implementation:
├── CSV Import: 2 недели
├── REST API Client: 2 недели
├── REST API Server: 2 недели
├── Delta sync: 1 неделя
├── Gradebook: 2 недели
├── Admin UI: 1 неделя
├── Testing & certification: 2 недели
└── Итого: 10-12 недель

IMS Certification (optional):
├── Conformance testing: 1-2 недели
├── Bug fixes: 1 неделя
├── Documentation: 1 неделя
└── Итого: 3-4 недели дополнительно
```

---

### 9.7 Источники для дальнейшего изучения

#### Официальные спецификации:

| Ресурс | Ссылка | Описание |
|--------|--------|----------|
| **OneRoster 1.1** | [imsglobal.org/oneroster](https://www.imsglobal.org/activity/onerosterlis) | Спецификация |
| **OneRoster 1.2** | [imsglobal.org](https://www.imsglobal.org/spec/oneroster/v1p2) | Новая версия |
| **REST Binding** | [imsglobal.org](https://www.imsglobal.org/oneroster-v11-final-specification) | API spec |
| **CSV Binding** | [imsglobal.org](https://www.imsglobal.org/oneroster-v11-final-csv-tables) | CSV format |

#### Tools & Libraries:

| Инструмент | Тип | Ссылка |
|------------|-----|--------|
| **OneRoster Validator** | Testing | IMS provides |
| **Clever** | Middleware | [clever.com](https://clever.com) |
| **ClassLink** | Middleware | [classlink.com](https://classlink.com) |

#### Примеры интеграций:

| Продукт | Тип | OneRoster |
|---------|-----|-----------|
| **Canvas LMS** | LMS | Full support |
| **Google Classroom** | LMS | Full support |
| **Schoology** | LMS | Full support |
| **PowerSchool** | SIS | Provider |
| **Infinite Campus** | SIS | Provider |

---

### 9.8 Чек-лист реализации

```
Phase 1: CSV Import (Week 1-2)
□ CSV parser for all entity types
□ Validation rules per OneRoster spec
□ ID mapping table
□ Import wizard UI
□ Error reporting
□ Rollback capability

Phase 2: Data Mapping (Week 3-4)
□ User role mapping configuration
□ Organization mapping
□ Course → internal course mapping
□ Class → internal class mapping
□ Enrollment sync logic
□ Deactivation handling

Phase 3: REST API Client (Week 5-6)
□ OAuth 2.0 client credentials flow
□ All GET endpoints
□ Pagination handling
□ Rate limiting
□ Error handling
□ Delta sync (filtering)

Phase 4: REST API Server (Week 7-8)
□ OAuth 2.0 token endpoint
□ All required GET endpoints
□ Pagination
□ Filtering (status, date)
□ HTTPS only

Phase 5: Gradebook (Week 9-10)
□ LineItem sync
□ Result sync
□ Bi-directional grades
□ Category mapping

Phase 6: Admin & Testing (Week 11-12)
□ Connection management UI
□ Sync scheduling
□ Monitoring dashboard
□ Sync history/logs
□ Integration tests
□ IMS conformance tests
```

---

*Продолжение документа с объяснением следующих функций будет добавлено в следующих разделах.*
