import React, { useState, useRef } from 'react';
import { 
  Bold, Italic, Heading1, Heading2, List, ListOrdered, 
  Quote, Code, Table as TableIcon, Link as LinkIcon, 
  Image as ImageIcon, Eye, EyeOff, FileText, CheckSquare, X
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Badge } from '@/components/ui/badge';
import { cn } from '@/lib/utils';

interface MarkdownEditorProps {
  value: string;
  onChange: (v: string) => void;
  className?: string;
}

const ToolbarButton = ({ icon: Icon, onClick, tooltip }: any) => (
  <button 
    type="button"
    onClick={onClick}
    className="p-2 text-slate-500 hover:text-indigo-600 hover:bg-white hover:shadow-sm rounded-lg transition-all"
    title={tooltip}
  >
    <Icon size={16} />
  </button>
);

export const MarkdownEditor: React.FC<MarkdownEditorProps> = ({ value, onChange, className }) => {
  const [isPreview, setIsPreview] = useState(false);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  const insertText = (textToInsert: string) => {
    const textarea = textareaRef.current;
    if (!textarea) {
      onChange(value + textToInsert);
      return;
    }

    const start = textarea.selectionStart;
    const end = textarea.selectionEnd;
    const text = textarea.value;
    const before = text.substring(0, start);
    const after = text.substring(end);
    
    onChange(before + textToInsert + after);
    
    setTimeout(() => {
      textarea.focus();
      textarea.setSelectionRange(start + textToInsert.length, start + textToInsert.length);
    }, 0);
  };

  const wrapText = (prefix: string, suffix: string = '') => {
    const textarea = textareaRef.current;
    if (!textarea) return;

    const start = textarea.selectionStart;
    const end = textarea.selectionEnd;
    const text = textarea.value;
    const before = text.substring(0, start);
    const selection = text.substring(start, end);
    const after = text.substring(end);

    const newText = before + prefix + selection + suffix + after;
    onChange(newText);
    
    setTimeout(() => {
      textarea.focus();
      textarea.setSelectionRange(start + prefix.length, end + prefix.length);
    }, 0);
  };

  const renderMarkdown = (markdown: string) => {
    if (!markdown) return <div className="flex flex-col items-center justify-center h-full text-slate-400"><FileText size={48} className="mb-4 opacity-20" /><p>Content preview will appear here.</p></div>;

    const lines = markdown.split('\n');
    const elements: React.ReactNode[] = [];
    
    lines.forEach((line, i) => {
      const trimmed = line.trim();
      if (trimmed === '') return;
      
      if (trimmed.startsWith('# ')) { elements.push(<h1 key={i} className="text-2xl font-bold mb-4 text-slate-900 border-b border-slate-100 pb-2">{trimmed.substring(2)}</h1>); return; }
      if (trimmed.startsWith('## ')) { elements.push(<h2 key={i} className="text-xl font-bold mb-3 mt-6 text-slate-800">{trimmed.substring(3)}</h2>); return; }
      if (trimmed.startsWith('### ')) { elements.push(<h3 key={i} className="text-lg font-bold mb-2 mt-4 text-slate-800">{trimmed.substring(4)}</h3>); return; }
      if (trimmed.startsWith('> ')) { elements.push(<blockquote key={i} className="border-l-4 border-indigo-500 pl-4 italic my-4 text-slate-600 bg-slate-50 py-2 pr-2 rounded-r">{trimmed.substring(2)}</blockquote>); return; }
      
      elements.push(<p key={i} className="mb-2 text-slate-600 leading-relaxed">{trimmed}</p>);
    });
    
    return <div className="prose prose-slate prose-sm max-w-none">{elements}</div>;
  };

  return (
    <div className={cn("flex flex-col border border-slate-200 rounded-xl overflow-hidden bg-white shadow-sm transition-all focus-within:ring-2 focus-within:ring-indigo-100 focus-within:border-indigo-300 relative", className)}>
      <div className="flex items-center gap-1 p-2 border-b border-slate-100 bg-slate-50 overflow-x-auto no-scrollbar">
        <ToolbarButton icon={Bold} onClick={() => wrapText('**', '**')} tooltip="Bold" />
        <ToolbarButton icon={Italic} onClick={() => wrapText('*', '*')} tooltip="Italic" />
        <div className="w-px h-4 bg-slate-300 mx-1 flex-shrink-0" />
        <ToolbarButton icon={Heading1} onClick={() => wrapText('# ')} tooltip="Heading 1" />
        <ToolbarButton icon={Heading2} onClick={() => wrapText('## ')} tooltip="Heading 2" />
        <div className="w-px h-4 bg-slate-300 mx-1 flex-shrink-0" />
        <ToolbarButton icon={List} onClick={() => wrapText('- ')} tooltip="Bullet List" />
        <ToolbarButton icon={ListOrdered} onClick={() => wrapText('1. ')} tooltip="Numbered List" />
        <div className="w-px h-4 bg-slate-300 mx-1 flex-shrink-0" />
        <ToolbarButton icon={Quote} onClick={() => wrapText('> ')} tooltip="Quote" />
        <ToolbarButton icon={Code} onClick={() => wrapText('`', '`')} tooltip="Code" />
        
        <div className="flex-1" />
        <button 
          onClick={() => setIsPreview(!isPreview)} 
          className={cn(
            "flex items-center gap-2 px-3 py-1.5 rounded-lg text-xs font-bold transition-colors ml-2 border flex-shrink-0",
            isPreview ? "bg-indigo-100 text-indigo-700 border-indigo-200" : "bg-white text-slate-600 border-slate-200 hover:bg-slate-50"
          )}
        >
          {isPreview ? <><EyeOff size={14} /> Edit</> : <><Eye size={14} /> Preview</>}
        </button>
      </div>

      <div className="relative min-h-[400px]">
        {isPreview ? (
          <div className="absolute inset-0 bg-white p-8 overflow-y-auto">
            {renderMarkdown(value)}
          </div>
        ) : (
          <textarea 
            ref={textareaRef}
            className="w-full h-full p-6 outline-none resize-none font-mono text-sm leading-relaxed text-slate-800 bg-white min-h-[400px]"
            placeholder="# Write your lesson content here..."
            value={value}
            onChange={(e) => onChange(e.target.value)}
            spellCheck={false}
          />
        )}
      </div>
      
      <div className="bg-slate-50 border-t border-slate-100 px-4 py-2 text-[10px] text-slate-400 flex justify-between font-mono items-center">
        <span>Markdown supported</span>
        <span>{value?.length || 0} chars</span>
      </div>
    </div>
  );
};
