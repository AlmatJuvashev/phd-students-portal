import * as React from "react";

export function DropdownMenu({
  trigger,
  children,
}: {
  trigger: React.ReactNode;
  children: React.ReactNode;
}) {
  const [open, setOpen] = React.useState(false);
  const dropdownRef = React.useRef<HTMLDivElement>(null);

  // Close dropdown when clicking outside
  React.useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (
        dropdownRef.current &&
        !dropdownRef.current.contains(event.target as Node)
      ) {
        setOpen(false);
      }
    };
    if (open) {
      document.addEventListener("mousedown", handleClickOutside);
    }
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, [open]);

  return (
    <div className="relative inline-block" ref={dropdownRef}>
      <div onClick={() => setOpen((o) => !o)}>{trigger}</div>
      {open && (
        <div className="absolute right-0 mt-2 min-w-[12rem] rounded-xl border border-border bg-card dark:bg-card shadow-lg z-50 p-1.5 animate-in fade-in-0 zoom-in-95 slide-in-from-top-2 duration-200">
          {children}
        </div>
      )}
    </div>
  );
}

export function DropdownItem({
  onClick,
  children,
}: {
  onClick?: () => void;
  children: React.ReactNode;
}) {
  return (
    <button
      onClick={onClick}
      className="block w-full text-left rounded-lg px-3 py-2.5 text-sm hover:bg-primary/10 dark:hover:bg-primary/20 transition-colors duration-150 focus:outline-none focus:bg-primary/10"
    >
      {children}
    </button>
  );
}
