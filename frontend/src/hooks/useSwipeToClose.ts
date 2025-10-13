import { useEffect, useRef } from "react";

interface UseSwipeToCloseOptions {
  onClose: () => void;
  enabled?: boolean;
  threshold?: number; // pixels to trigger close
}

/**
 * Hook to enable swipe-down-to-close gesture for bottom sheets on mobile
 */
export function useSwipeToClose({
  onClose,
  enabled = true,
  threshold = 100,
}: UseSwipeToCloseOptions) {
  const startYRef = useRef<number | null>(null);
  const currentYRef = useRef<number | null>(null);
  const isDraggingRef = useRef(false);

  useEffect(() => {
    if (!enabled) return;

    const handleTouchStart = (e: TouchEvent) => {
      const target = e.target as HTMLElement;
      
      // Only start tracking if touch starts on the drag handle or header
      const isDragHandle = target.closest('[data-drag-handle]');
      const isHeader = target.closest('[data-sheet-header]');
      
      if (isDragHandle || isHeader) {
        startYRef.current = e.touches[0].clientY;
        currentYRef.current = e.touches[0].clientY;
        isDraggingRef.current = true;
      }
    };

    const handleTouchMove = (e: TouchEvent) => {
      if (!isDraggingRef.current || startYRef.current === null) return;

      currentYRef.current = e.touches[0].clientY;
      const deltaY = currentYRef.current - startYRef.current;

      // Only allow downward swipes (positive deltaY)
      if (deltaY > 0) {
        // Optionally prevent scroll while dragging
        const target = e.target as HTMLElement;
        const scrollable = target.closest('[data-sheet-content]');
        
        if (scrollable && scrollable.scrollTop === 0) {
          e.preventDefault();
        }
      }
    };

    const handleTouchEnd = () => {
      if (!isDraggingRef.current || startYRef.current === null || currentYRef.current === null) {
        isDraggingRef.current = false;
        return;
      }

      const deltaY = currentYRef.current - startYRef.current;

      // Trigger close if swipe down exceeds threshold
      if (deltaY > threshold) {
        onClose();
      }

      // Reset
      startYRef.current = null;
      currentYRef.current = null;
      isDraggingRef.current = false;
    };

    document.addEventListener("touchstart", handleTouchStart, { passive: true });
    document.addEventListener("touchmove", handleTouchMove, { passive: false });
    document.addEventListener("touchend", handleTouchEnd, { passive: true });
    document.addEventListener("touchcancel", handleTouchEnd, { passive: true });

    return () => {
      document.removeEventListener("touchstart", handleTouchStart);
      document.removeEventListener("touchmove", handleTouchMove);
      document.removeEventListener("touchend", handleTouchEnd);
      document.removeEventListener("touchcancel", handleTouchEnd);
    };
  }, [enabled, threshold, onClose]);
}
