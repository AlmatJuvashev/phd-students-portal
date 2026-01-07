import { useEditor, EditorContent } from '@tiptap/react';
import StarterKit from '@tiptap/starter-kit';
import Link from '@tiptap/extension-link';
import Image from '@tiptap/extension-image';
import Underline from '@tiptap/extension-underline';
import { Button } from './button';
import { 
  Bold, Italic, Strikethrough, List, ListOrdered, Link as LinkIcon, 
  Image as ImageIcon, Undo, Redo, Underline as UnderlineIcon, Quote 
} from 'lucide-react';
import { cn } from '@/lib/utils';
import { useEffect } from 'react';

interface RichTextEditorProps {
  value: string;
  onChange: (value: string) => void;
  placeholder?: string;
  className?: string;
}

export const RichTextEditor = ({ value, onChange, placeholder, className }: RichTextEditorProps) => {
  const editor = useEditor({
    extensions: [
      StarterKit,
      Underline,
      Link.configure({
        openOnClick: false,
      }),
      Image,
    ],
    content: value,
    editorProps: {
      attributes: {
        class: cn(
          "min-h-[150px] w-full rounded-md border border-input bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-1 focus-visible:ring-ring disabled:cursor-not-allowed disabled:opacity-50 overflow-auto",
           className
        ),
      },
    },
    onUpdate: ({ editor }) => {
      onChange(editor.getHTML());
    },
  });

  // Sync value if changed externally
  useEffect(() => {
    if (editor && value !== editor.getHTML()) {
      // Only set content if it's actually different to avoid cursor jumping
      // Simple check, real world might need deeper diff or managing focus
      if (editor.getText() === '' && value === '') return;
      // editor.commands.setContent(value); 
      // NOTE: setting content on every prop change causes cursor reset issues.
      // Usually better to let editor manage state and only sync on mount or specific reset.
      // For now, we assume uncontrolled or careful controlled usage.
      // If we strictly need controlled, we check if content is same.
    }
  }, [value, editor]);

  if (!editor) {
    return null;
  }

  const toggleBold = () => editor.chain().focus().toggleBold().run();
  const toggleItalic = () => editor.chain().focus().toggleItalic().run();
  const toggleUnderline = () => editor.chain().focus().toggleUnderline().run();
  const toggleStrike = () => editor.chain().focus().toggleStrike().run();
  const toggleBulletList = () => editor.chain().focus().toggleBulletList().run();
  const toggleOrderedList = () => editor.chain().focus().toggleOrderedList().run();
  const toggleBlockquote = () => editor.chain().focus().toggleBlockquote().run();
  
  const addLink = () => {
    const previousUrl = editor.getAttributes('link').href;
    const url = window.prompt('URL', previousUrl);
    
    if (url === null) return;
    if (url === '') {
      editor.chain().focus().extendMarkRange('link').unsetLink().run();
      return;
    }
    editor.chain().focus().extendMarkRange('link').setLink({ href: url }).run();
  };

  const addImage = () => {
    const url = window.prompt('Image URL');
    if (url) {
      editor.chain().focus().setImage({ src: url }).run();
    }
  };

  return (
    <div className="flex flex-col gap-2">
      <div className="flex flex-wrap gap-1 p-1 border rounded-md bg-muted/20">
        <ToolbarButton onClick={toggleBold} isActive={editor.isActive('bold')} icon={<Bold size={16} />} />
        <ToolbarButton onClick={toggleItalic} isActive={editor.isActive('italic')} icon={<Italic size={16} />} />
        <ToolbarButton onClick={toggleUnderline} isActive={editor.isActive('underline')} icon={<UnderlineIcon size={16} />} />
        <ToolbarButton onClick={toggleStrike} isActive={editor.isActive('strike')} icon={<Strikethrough size={16} />} />
        
        <div className="w-px bg-border mx-1 h-6 self-center" />
        
        <ToolbarButton onClick={toggleBulletList} isActive={editor.isActive('bulletList')} icon={<List size={16} />} />
        <ToolbarButton onClick={toggleOrderedList} isActive={editor.isActive('orderedList')} icon={<ListOrdered size={16} />} />
        <ToolbarButton onClick={toggleBlockquote} isActive={editor.isActive('blockquote')} icon={<Quote size={16} />} />

        <div className="w-px bg-border mx-1 h-6 self-center" />

        <ToolbarButton onClick={addLink} isActive={editor.isActive('link')} icon={<LinkIcon size={16} />} />
        <ToolbarButton onClick={addImage} isActive={false} icon={<ImageIcon size={16} />} />

        <div className="w-px bg-border mx-1 h-6 self-center" />

        <ToolbarButton onClick={() => editor.chain().focus().undo().run()} isActive={false} icon={<Undo size={16} />} />
        <ToolbarButton onClick={() => editor.chain().focus().redo().run()} isActive={false} icon={<Redo size={16} />} />
      </div>
      
      <EditorContent editor={editor} />
    </div>
  );
};

const ToolbarButton = ({ onClick, isActive, icon }: { onClick: () => void, isActive: boolean, icon: React.ReactNode }) => (
  <Button
    type="button"
    variant={isActive ? "secondary" : "ghost"}
    size="icon"
    className="h-8 w-8"
    onClick={onClick}
  >
    {icon}
  </Button>
);
