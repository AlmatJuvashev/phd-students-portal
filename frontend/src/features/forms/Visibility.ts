export function evalVisible(values: Record<string, any>, expr?: string): boolean {
  if (!expr) return true;
  try {
    const mEq = expr.match(/form\.([a-zA-Z0-9_]+)\s*==\s*(true|false)/);
    if (mEq) {
      const key = mEq[1];
      const val = mEq[2] === "true";
      if (!Object.prototype.hasOwnProperty.call(values, key)) return false;
      return !!values[key] === val;
    }
    const mNeq = expr.match(/form\.([a-zA-Z0-9_]+)\s*!=\s*(true|false)/);
    if (mNeq) {
      const key = mNeq[1];
      const val = mNeq[2] === "true";
      if (!Object.prototype.hasOwnProperty.call(values, key)) return false;
      return !!values[key] !== val;
    }
    if (expr.includes("&&") || expr.includes("||")) {
      const replaced = expr.replace(/form\.([a-zA-Z0-9_]+)/g, (s, k) => {
        return Object.prototype.hasOwnProperty.call(values, k)
          ? JSON.stringify(!!values[k])
          : "undefined";
      });
      // eslint-disable-next-line no-new-func
      return Function(`return (${replaced});`)();
    }
    return true;
  } catch (e) {
    console.error("Error evaluating visibility expression:", expr, e);
    return true;
  }
}

