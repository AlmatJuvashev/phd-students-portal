import { api } from "./client";

export async function uploadReviewedDocument(
  attachmentId: string,
  payload: {
    document_version_id: string;
  }
) {
  return api(`/admin/attachments/${attachmentId}/reviewed-document`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export async function presignReviewedDocumentUpload(
  attachmentId: string,
  payload: {
    filename: string;
    content_type: string;
    size_bytes: number;
  }
) {
  return api(`/admin/attachments/${attachmentId}/presign`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
}

export async function attachReviewedDocument(
  attachmentId: string,
  payload: {
    object_key: string;
    filename: string;
    content_type: string;
    size_bytes: number;
    etag?: string;
  }
) {
  return api(`/admin/attachments/${attachmentId}/attach-reviewed`, {
    method: "POST",
    body: JSON.stringify(payload),
  });
}
