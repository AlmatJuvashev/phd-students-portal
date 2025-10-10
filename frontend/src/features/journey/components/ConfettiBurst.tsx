import { useEffect, useState } from "react";

export function ConfettiBurst({ trigger }: { trigger: boolean }) {
  const [show, setShow] = useState(false);
  useEffect(() => {
    if (trigger) {
      setShow(true);
      const t = setTimeout(() => setShow(false), 1200);
      return () => clearTimeout(t);
    }
  }, [trigger]);
  if (!show) return null;
  return (
    <div className="pointer-events-none fixed inset-0 z-[80] overflow-hidden">
      {Array.from({ length: 80 }).map((_, i) => (
        <span
          key={i}
          className="absolute block w-1.5 h-3 rounded-sm"
          style={{
            left: Math.random() * 100 + "%",
            top: "-10px",
            backgroundColor: randomColor(),
            transform: `rotate(${Math.random() * 360}deg)`,
            animation: `fall ${800 + Math.random() * 800}ms ease-out forwards` as any,
          }}
        />
      ))}
      <style>
        {`
          @keyframes fall { to { transform: translateY(110vh) rotate(720deg); opacity: 0.8; } }
        `}
      </style>
    </div>
  );
}

function randomColor() {
  const colors = ["#ef4444", "#f59e0b", "#10b981", "#3b82f6", "#8b5cf6", "#14b8a6"];
  return colors[Math.floor(Math.random() * colors.length)];
}

