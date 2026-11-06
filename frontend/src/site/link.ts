// Situation links (ADR 026, 7): the payload JSON, deflate-raw, then
// base64url with no padding. The browser has CompressionStream, so no
// library; Go's compress/flate reads the same bytes (11.3).

// linkWarnAt: above this many characters some chats cut a link.
export const linkWarnAt = 2000;

// compact gives the payload with no spaces when it is JSON, so links
// stay short; text that is no JSON stays as it is.
export function compact(payload: string): string {
  try {
    return JSON.stringify(JSON.parse(payload));
  } catch {
    return payload;
  }
}

async function pipe(bytes: Uint8Array<ArrayBuffer>, through: CompressionStream | DecompressionStream): Promise<Uint8Array> {
  const stream = new Blob([bytes]).stream().pipeThrough(through);
  return new Uint8Array(await new Response(stream).arrayBuffer());
}

// toBase64url is base64 with - and _, and no = at the end.
export function toBase64url(bytes: Uint8Array): string {
  let bin = '';
  for (const b of bytes) {
    bin += String.fromCharCode(b);
  }
  return btoa(bin).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
}

// fromBase64url reverses toBase64url; it throws on other characters.
export function fromBase64url(text: string): Uint8Array<ArrayBuffer> {
  if (!/^[A-Za-z0-9_-]*$/.test(text)) {
    throw new Error('not base64url');
  }
  const bin = atob(text.replace(/-/g, '+').replace(/_/g, '/'));
  return Uint8Array.from(bin, (c) => c.charCodeAt(0));
}

// encodePayload is the <data> of a check link.
export async function encodePayload(payload: string): Promise<string> {
  return toBase64url(await pipe(new TextEncoder().encode(payload), new CompressionStream('deflate-raw')));
}

// decodePayload is the payload of a check link's <data>; it throws
// when the data is cut or broken.
export async function decodePayload(data: string): Promise<string> {
  const bytes = await pipe(fromBase64url(data), new DecompressionStream('deflate-raw'));
  return new TextDecoder('utf-8', { fatal: true }).decode(bytes);
}
