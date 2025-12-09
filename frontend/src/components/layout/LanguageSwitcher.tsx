import React from "react";
import { useTranslation } from "react-i18next";
import { Globe } from "lucide-react";
import { Button } from "@/components/ui/button";
import { DropdownMenu, DropdownItem } from "@/components/ui/dropdown-menu";

export function LanguageSwitcher() {
  const { i18n } = useTranslation("common");

  return (
    <DropdownMenu
      trigger={
        <Button
          variant="ghost"
          size="sm"
          className="gap-2 h-9 px-2 text-muted-foreground hover:text-foreground"
        >
          <Globe className="h-4 w-4" />
          <span className="hidden sm:inline">
            {i18n.language.toUpperCase()}
          </span>
        </Button>
      }
    >
      <DropdownItem onClick={() => i18n.changeLanguage("ru")}>
        <div
          className={`flex items-center gap-3 ${
            i18n.language === "ru" ? "font-semibold text-primary" : ""
          }`}
        >
          <span className="text-lg">🇷🇺</span>
          <span>Русский</span>
        </div>
      </DropdownItem>
      <DropdownItem onClick={() => i18n.changeLanguage("kz")}>
        <div
          className={`flex items-center gap-3 ${
            i18n.language === "kz" ? "font-semibold text-primary" : ""
          }`}
        >
          <span className="text-lg">🇰🇿</span>
          <span>Қазақша</span>
        </div>
      </DropdownItem>
      <DropdownItem onClick={() => i18n.changeLanguage("en")}>
        <div
          className={`flex items-center gap-3 ${
            i18n.language === "en" ? "font-semibold text-primary" : ""
          }`}
        >
          <span className="text-lg">🇬🇧</span>
          <span>English</span>
        </div>
      </DropdownItem>
    </DropdownMenu>
  );
}
