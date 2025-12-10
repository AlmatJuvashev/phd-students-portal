import React from 'react';
import { motion } from 'framer-motion';
import { cn } from '../lib/utils';

interface JourneyPathProps {
  state: 'done' | 'active' | 'locked' | 'waiting' | 'needs_fixes' | 'under_review';
  isLast: boolean;
}

export const JourneyPath: React.FC<JourneyPathProps> = ({ state, isLast }) => {
  if (isLast) return null;

  return (
    <div className="absolute left-[2.1rem] sm:left-[2.85rem] top-14 h-[calc(100%-1rem)] w-1.5 flex items-start justify-center -z-10">
      {/* Background Track (The "Unexplored" Road) */}
      <div className="h-full w-full flex flex-col items-center overflow-hidden opacity-40">
        {/* We use a repeating dashed border for the locked state to look like a treasure map path */}
        <div className="h-full w-0 border-l-[3px] border-dashed border-slate-300" />
      </div>

      {/* Progress Fill (The "Traveled" Road - DONE state) */}
      <motion.div
        initial={{ height: "0%" }}
        animate={{ 
          height: state === 'done' ? '100%' : '0%'
        }}
        transition={{ duration: 0.8, ease: "easeInOut" }}
        className={cn(
          "absolute top-0 w-full rounded-full bg-emerald-500 origin-top overflow-hidden",
          state === 'done' && "shadow-[0_0_10px_rgba(16,185,129,0.6)]"
        )}
      >
        {/* 'Done' State: Occasional shine traveling down the wire */}
        {state === 'done' && (
          <>
            <motion.div
              className="absolute left-0 w-full h-20 bg-gradient-to-b from-transparent via-white/60 to-transparent"
              animate={{ top: ["-20%", "120%"] }}
              transition={{ 
                duration: 3, 
                repeat: Infinity, 
                ease: "easeInOut",
                repeatDelay: 2
              }}
            />
            {/* Subtle pulse of the entire line to show it is active/powered */}
            <motion.div 
               className="absolute inset-0 bg-emerald-400/30"
               animate={{ opacity: [0, 0.3, 0] }}
               transition={{ duration: 3, repeat: Infinity, ease: "easeInOut" }}
            />
          </>
        )}
      </motion.div>
      
      {/* Active Gradient (The "Current" Path) */}
      {state === 'active' && (
        <motion.div
          initial={{ height: "0%" }}
          animate={{ height: "100%" }}
          className="absolute top-0 w-full h-full"
        >
          {/* The main gradient line */}
          <div className="h-full w-full rounded-full bg-gradient-to-b from-primary-500 via-primary-300 to-transparent opacity-90" />
          
          {/* Pulsing Aura Overlay - Breathing effect */}
          <motion.div
            className="absolute inset-0 rounded-full bg-primary-400 blur-md"
            animate={{ 
              opacity: [0.3, 0.6, 0.3],
              scaleX: [1, 1.5, 1] 
            }}
            transition={{ duration: 2, repeat: Infinity, ease: "easeInOut" }}
          />

          {/* Continuous Glowing Particles Stream */}
          <div className="absolute inset-0 overflow-hidden">
             {Array.from({ length: 8 }).map((_, i) => (
                <motion.div
                  key={i}
                  className="absolute bg-white rounded-full shadow-[0_0_8px_rgba(255,255,255,1)] z-10"
                  style={{
                    width: '3px',
                    height: '3px',
                    left: '50%',
                    marginLeft: '-1.5px' // Perfectly centered
                  }}
                  initial={{ top: "-10%", opacity: 0 }}
                  animate={{
                    top: "110%",
                    opacity: [0, 1, 1, 0] // Fade in, glow, fade out
                  }}
                  transition={{
                    duration: 3,
                    repeat: Infinity,
                    ease: "linear",
                    delay: i * 0.4, // Consistent flow spacing
                    repeatDelay: 0
                  }}
                />
             ))}
          </div>
        </motion.div>
      )}
      
      {/* Subtle Energy Mote (Done State) - Rare particle to show it's still alive */}
      {state === 'done' && (
        <motion.div
          className="absolute bg-emerald-100 rounded-full blur-[2px] shadow-lg"
          style={{ width: '4px', height: '4px', left: '1px' }}
          animate={{ top: ['0%', '100%'], opacity: [0, 0.8, 0] }}
          transition={{ duration: 5, repeat: Infinity, ease: "linear", repeatDelay: 1.5 }}
        />
      )}
    </div>
  );
};
