import * as React from "react";
import { format } from "date-fns";
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
  const [date, setDate] = React.useState<Date | undefined>(
    value ? new Date(value) : undefined
  );

  const locale = localeMap[i18n.language as keyof typeof localeMap] || enUS;

  React.useEffect(() => {
    if (value) {
      setDate(new Date(value));
    } else {
      setDate(undefined);
    }
  }, [value]);

  const handleSelect = (selectedDate: Date | undefined) => {
    setDate(selectedDate);
    if (selectedDate) {
      // Format as YYYY-MM-DD for compatibility with backend
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
      <PopoverContent className="w-auto p-0">
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
