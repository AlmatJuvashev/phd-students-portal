import * as React from "react";
import { format, parse } from "date-fns";
import { enUS, ru, kk } from "date-fns/locale";
import { Calendar as CalendarIcon } from "lucide-react";
import { useTranslation } from "react-i18next";

import { cn } from "@/lib/utils";
import { Button } from "@/components/ui/button";
import { Calendar } from "@/components/ui/calendar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover";

interface DatePickerProps {
  value?: string;
  onChange?: (value: string) => void;
  disabled?: boolean;
  placeholder?: string;
  className?: string;
}

const localeMap = {
  en: enUS,
  ru: ru,
  kz: kk,
};

export function DatePicker({
  value,
  onChange,
  disabled,
  placeholder = "Pick a date",
  className,
}: DatePickerProps) {
  const { i18n } = useTranslation();
  const [date, setDate] = React.useState<Date | undefined>(() => {
    if (!value) return undefined;
    // Parse ISO date string (YYYY-MM-DD) as local date to avoid timezone shifts
    const parsed = parse(value, "yyyy-MM-dd", new Date());
    // Fallback if parsing fails (e.g. if value is full ISO string)
    return isNaN(parsed.getTime()) ? new Date(value) : parsed;
  });

  const locale = localeMap[i18n.language as keyof typeof localeMap] || enUS;

  React.useEffect(() => {
    if (value) {
      const parsed = parse(value, "yyyy-MM-dd", new Date());
      setDate(isNaN(parsed.getTime()) ? new Date(value) : parsed);
    } else {
      setDate(undefined);
    }
  }, [value]);

  const handleSelect = (selectedDate: Date | undefined) => {
    // When a date is selected from Calendar (local time), format it back to YYYY-MM-DD
    setDate(selectedDate);
    if (selectedDate) {
      const formatted = format(selectedDate, "yyyy-MM-dd");
      onChange?.(formatted);
    } else {
      onChange?.("");
    }
  };

  return (
    <Popover>
      <PopoverTrigger asChild>
        <Button
          type="button" 
          variant={"outline"}
          className={cn(
            "w-full justify-start text-left font-normal",
            !date && "text-muted-foreground",
            className
          )}
          disabled={disabled}
        >
          <CalendarIcon className="mr-2 h-4 w-4" />
          {date ? format(date, "PPP", { locale }) : <span>{placeholder}</span>}
        </Button>
      </PopoverTrigger>
      <PopoverContent className="w-auto p-0 z-[100]">
        <Calendar
          mode="single"
          selected={date}
          onSelect={handleSelect}
          initialFocus
          locale={locale}
        />
      </PopoverContent>
    </Popover>
  );
}
