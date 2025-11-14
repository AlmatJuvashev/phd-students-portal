export function NotFoundIllustration({ className = "w-64 h-64" }: { className?: string }) {
  // Simple friendly 404 illustration using inline SVG and currentColor accents
  return (
    <svg
      viewBox="0 0 320 220"
      role="img"
      aria-label="Not found illustration"
      className={className}
    >
      <defs>
        <linearGradient id="nfGrad" x1="0" x2="1" y1="0" y2="1">
          <stop offset="0%" stopColor="hsl(210 60% 60%)" />
          <stop offset="100%" stopColor="hsl(260 60% 60%)" />
        </linearGradient>
        <filter id="shadow" x="-20%" y="-20%" width="140%" height="140%">
          <feDropShadow dx="0" dy="6" stdDeviation="8" floodOpacity="0.15" />
        </filter>
      </defs>
      <rect x="0" y="0" width="320" height="220" rx="16" fill="url(#nfGrad)" opacity="0.08" />
      <g filter="url(#shadow)">
        <circle cx="90" cy="110" r="42" fill="currentColor" opacity="0.08" />
        <circle cx="230" cy="110" r="42" fill="currentColor" opacity="0.08" />
        <rect x="65" y="70" width="50" height="80" rx="12" fill="white" stroke="currentColor" strokeOpacity="0.15" />
        <rect x="205" y="70" width="50" height="80" rx="12" fill="white" stroke="currentColor" strokeOpacity="0.15" />
        <rect x="135" y="58" width="50" height="104" rx="12" fill="white" stroke="currentColor" strokeOpacity="0.15" />
        <circle cx="160" cy="96" r="10" fill="currentColor" opacity="0.25" />
        <rect x="148" y="112" width="24" height="4" rx="2" fill="currentColor" opacity="0.25" />
      </g>
      <g>
        <text x="50%" y="188" textAnchor="middle" fontFamily="system-ui, -apple-system, Segoe UI, Roboto" fontSize="28" fill="currentColor" fillOpacity="0.3">
          404
        </text>
      </g>
    </svg>
  );
}

