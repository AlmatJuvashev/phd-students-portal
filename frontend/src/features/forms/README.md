# Unified Form System

**Phase 1.4 - Form Unification**

A comprehensive, reusable form management system that eliminates code duplication across form components.

## 🎯 Purpose

Before unification, form logic was duplicated across:
- `ChecklistDetails.tsx` - Checklist forms
- `FormTaskDetails.tsx` - Complex forms with multiple field types
- Custom scene components - Each with its own validation

This system provides:
- ✅ Unified state management
- ✅ Automatic validation
- ✅ Consistent submission logic
- ✅ Reusable UI components
- ✅ TypeScript type safety

## 📦 Components

### FormProvider

Context provider that manages form state and logic.

```tsx
<FormProvider
  node={node}
  initial={initialValues}
  onSubmit={handleSubmit}
  canEdit={true}
  disabled={false}
>
  {/* Your form content */}
</FormProvider>
```

**Props:**
- `node: NodeVM` - Node definition with fields
- `initial?: Record<string, any>` - Initial values
- `onSubmit?: (payload: any) => void` - Submit handler
- `canEdit?: boolean` - Override edit permissions
- `disabled?: boolean` - Disable all interactions

### useFormContext()

Hook to access form state from any child component.

```tsx
const {
  values,           // Current form values
  setField,         // Update single field
  setValues,        // Update multiple fields
  isValid,          // Form validation status
  canEdit,          // Edit permission
  readOnly,         // Computed readonly state
  submit,           // Submit handler
  saveDraft,        // Draft handler
  evalVisible,      // Check field visibility
  getNextOnComplete // Calculate next node
} = useFormContext();
```

### FormFields

Automatically renders all form fields based on node configuration.

```tsx
<FormFields 
  types={["boolean"]}  // Optional: filter by field type
  spacing="normal"     // "tight" | "normal" | "loose"
  className=""
/>
```

**Supported field types:**
- `boolean` - Rendered as ChecklistItem
- `text` - Text input
- `select` - Dropdown with optional "other" field
- `date` - Date picker

### FormActions

Standard submit and draft buttons with optional confirmation modal.

```tsx
<FormActions
  showConfirm={true}      // Show confirmation modal
  submitLabel="Submit"    // Custom submit label
  draftLabel="Save"       // Custom draft label
  hideSubmit={false}      // Hide submit button
  hideDraft={false}       // Hide draft button
/>
```

Features:
- Auto-disabled when form invalid
- Touch-optimized (44px min height)
- Loading states
- Confirmation modal integration

### FormStatus

Displays submission status and timestamp.

```tsx
<FormStatus className="" />
```

Shows:
- "Form submitted (date: ...)" when readonly
- Hidden when form is editable
- Localized dates (ru/kz/en)

## 🔥 Usage Examples

### Simple Checklist Form

```tsx
import { FormProvider, FormFields, FormActions, FormStatus } from "@/features/forms";

function MyChecklist({ node, initial, onSubmit }) {
  return (
    <FormProvider node={node} initial={initial} onSubmit={onSubmit}>
      <FormFields types={["boolean"]} />
      <FormActions showConfirm={true} />
      <FormStatus />
    </FormProvider>
  );
}
```

### Complex Mixed Form

```tsx
function MyComplexForm({ node, initial, onSubmit }) {
  return (
    <FormProvider node={node} initial={initial} onSubmit={onSubmit}>
      <div className="space-y-4">
        <FormFields types={["text", "select"]} spacing="loose" />
        <FormFields types={["boolean"]} spacing="normal" />
        <FormActions hideSubmit={false} hideDraft={true} />
      </div>
    </FormProvider>
  );
}
```

### Custom Actions with Context

```tsx
function CustomForm({ node, initial, onSubmit }) {
  return (
    <FormProvider node={node} initial={initial} onSubmit={onSubmit}>
      <FormFields />
      <CustomActions />
    </FormProvider>
  );
}

function CustomActions() {
  const { isValid, submit, values } = useFormContext();
  
  return (
    <Button 
      onClick={() => submit({ customField: "value" })}
      disabled={!isValid || !values.someRequirement}
    >
      Custom Submit
    </Button>
  );
}
```

## 🎨 Features

### Automatic Validation

```typescript
// FormProvider automatically validates:
- Required fields are filled
- Hidden fields are excluded (via visible expression)
- Boolean fields are checked
- Text/select/date fields are not empty
```

### Next Node Resolution

```typescript
// Automatically resolves next node based on outcomes
getNextOnComplete() => {
  // 1. Check outcomes with "when" conditions
  // 2. Match current form state to conditions
  // 3. Return matching outcome's next node
  // 4. Fallback to first outcome or node.next
}
```

### Dependent Fields

```typescript
// Built-in clearing of dependent fields:
setField("hearing_happened", false)
  → auto-clears: remarks_exist, plan_prepared, remarks_resolved

setField("remarks_exist", false)
  → auto-clears: plan_prepared, remarks_resolved
```

## 🔄 Migration Guide

### Before (ChecklistDetails.tsx)

```tsx
const [values, setValues] = useState(initial);
const fields = node.requirements?.fields ?? [];
const ready = fields.every(f => !!values[f.key]);
const readOnly = node.state === "submitted" || ...;

<ChecklistItem
  checked={!!values[field.key]}
  onChange={(checked) => setValues(s => ({ ...s, [field.key]: checked }))}
/>

<Button onClick={handleSubmit} disabled={!ready}>
  Submit
</Button>
```

### After (Unified)

```tsx
<FormProvider node={node} initial={initial} onSubmit={handleSubmit}>
  <FormFields types={["boolean"]} />
  <FormActions showConfirm={true} />
</FormProvider>
```

**Benefits:**
- ❌ ~100 lines of code removed
- ✅ No manual state management
- ✅ No manual validation
- ✅ No duplicate logic
- ✅ Consistent behavior

## 📊 Metrics

**Before Phase 1.4:**
- ChecklistDetails: 173 lines
- FormTaskDetails: 208 lines
- Custom scenes: 150+ lines each
- Total duplication: ~500+ lines

**After Phase 1.4:**
- FormProvider: 243 lines (reusable)
- FormFields: 80 lines (reusable)
- FormActions: 75 lines (reusable)
- ChecklistDetailsUnified: 67 lines (−61% reduction)
- Estimated savings: 300+ lines across all forms

## 🚀 Next Steps

1. **Migrate ChecklistDetails** → Use unified system
2. **Migrate FormTaskDetails** → Simplify with FormFields
3. **Migrate custom scenes** → Replace with unified components
4. **Add tests** → Test FormProvider logic
5. **Document patterns** → Create style guide

## ⚙️ Technical Details

**State Management:**
- `useReducer` for form state (predictable updates)
- Smart field dependencies (auto-clear related fields)
- Memoized computed values (isValid, nextNode)

**Performance:**
- Memoized context value (prevent unnecessary re-renders)
- Lazy field rendering (hidden fields not mounted)
- Optimized validation (only required fields checked)

**Accessibility:**
- 44px touch targets on all buttons
- `aria-busy` states during submission
- Screen reader labels
- Keyboard navigation support

## 📝 API Reference

See inline TypeScript documentation in component files for full API details.
