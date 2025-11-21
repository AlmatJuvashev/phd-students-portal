import { X } from "lucide-react";
import { Button } from "./button";
import { useEffect, useState } from "react";

type DocumentPreviewDrawerProps = {
  isOpen: boolean;
  onClose: () => void;
  fileUrl: string;
  filename: string;
  contentType?: string;
};

export function DocumentPreviewDrawer({
  isOpen,
  onClose,
  fileUrl,
  filename,
  contentType,
}: DocumentPreviewDrawerProps) {
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!isOpen || !fileUrl) {
      setPreviewUrl(null);
      setError(null);
      return;
    }

    setError(null);
    setPreviewUrl(fileUrl);
  }, [isOpen, fileUrl]);

  if (!isOpen) return null;

  const getFileExtension = (name: string) => {
    const parts = name.split(".");
    return parts.length > 1 ? parts[parts.length - 1].toLowerCase() : "";
  };

  const ext = getFileExtension(filename);
  const isPDF = ext === "pdf" || contentType?.includes("pdf");
  const isImage =
    ["jpg", "jpeg", "png", "gif", "webp", "svg"].includes(ext) ||
    contentType?.startsWith("image/");
  const isDOCX =
    ext === "docx" ||
    contentType?.includes(
      "vnd.openxmlformats-officedocument.wordprocessingml.document"
    );

  const renderPreview = () => {
    if (error) {
      return (
        <div className="flex items-center justify-center h-full">
          <div className="text-center p-6">
            <p className="text-destructive mb-2">Failed to load preview</p>
            <p className="text-sm text-muted-foreground">{error}</p>
            <Button variant="outline" className="mt-4" onClick={onClose}>
              Close
            </Button>
          </div>
        </div>
      );
    }

    if (!previewUrl) {
      return (
        <div className="flex items-center justify-center h-full">
          <p className="text-muted-foreground">Loading preview...</p>
        </div>
      );
    }

    if (isPDF) {
      return (
        <iframe
          src={previewUrl}
          className="w-full h-full border-0"
          title={filename}
          onError={() => setError("Unable to display PDF file")}
        />
      );
    }

    if (isImage) {
      return (
        <div className="flex items-center justify-center h-full overflow-auto p-4">
          <img
            src={previewUrl}
            alt={filename}
            className="max-w-full max-h-full object-contain"
            onError={() => setError("Unable to display image")}
          />
        </div>
      );
    }

    if (isDOCX) {
      // Note: Google Docs Viewer works only with publicly accessible URLs
      // For localhost/private URLs, this will show an error
      const googleDocsUrl = `https://docs.google.com/viewer?url=${encodeURIComponent(
        previewUrl
      )}&embedded=true`;

      return (
        <div className="flex flex-col h-full">
          <div className="bg-amber-50 border-b border-amber-200 p-3">
            <p className="text-sm text-amber-900">
              ⚠️ DOCX preview requires a publicly accessible URL. Local files
              cannot be previewed. Please download to view.
            </p>
          </div>
          <iframe
            src={googleDocsUrl}
            className="w-full flex-1 border-0"
            title={filename}
            onError={() =>
              setError(
                "Unable to preview DOCX file. Google Docs Viewer requires public URLs."
              )
            }
          />
        </div>
      );
    }

    // Unsupported file type
    return (
      <div className="flex items-center justify-center h-full">
        <div className="text-center p-6">
          <p className="text-muted-foreground mb-2">
            Preview not available for this file type
          </p>
          <p className="text-sm text-muted-foreground mb-4">File: {filename}</p>
          <Button variant="outline" onClick={onClose}>
            Close
          </Button>
        </div>
      </div>
    );
  };

  return (
    <>
      {/* Overlay */}
      <div
        className="fixed inset-0 bg-black/50 z-40 transition-opacity"
        onClick={onClose}
      />

      {/* Drawer */}
      <div className="fixed inset-y-0 right-0 w-full md:w-3/4 lg:w-2/3 xl:w-1/2 bg-background shadow-xl z-50 flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between px-6 py-4 border-b">
          <div className="flex-1 min-w-0">
            <h2 className="text-lg font-semibold truncate">{filename}</h2>
            <p className="text-xs text-muted-foreground mt-0.5">
              {isPDF && "PDF Document"}
              {isImage && "Image Preview"}
              {isDOCX && "Word Document"}
              {!isPDF && !isImage && !isDOCX && "File Preview"}
            </p>
          </div>
          <Button
            variant="ghost"
            size="icon"
            onClick={onClose}
            className="flex-shrink-0 ml-4"
          >
            <X className="h-5 w-5" />
          </Button>
        </div>

        {/* Preview Content */}
        <div className="flex-1 overflow-hidden">{renderPreview()}</div>
      </div>
    </>
  );
}
