import { Link, useLocation } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/button";
import { NotFoundIllustration } from "./NotFoundIllustration";

export default function NotFound() {
  const { t } = useTranslation("common");
  const loc = useLocation();
  const isAdmin = loc.pathname.startsWith("/admin");
  const homeHref = isAdmin ? "/admin" : "/";

  return (
    <div className="min-h-[60vh] flex items-center justify-center px-4">
      <div className="max-w-xl text-center space-y-6">
        <div className="flex items-center justify-center text-primary">
          <NotFoundIllustration className="w-60 h-40 sm:w-72 sm:h-48" />
        </div>
        <h1 className="text-2xl sm:text-3xl font-bold">
          {t("errors.not_found_title", { defaultValue: "Page not found" })}
        </h1>
        <p className="text-muted-foreground">
          {t("errors.not_found_desc", {
            defaultValue:
              "The page you’re looking for doesn’t exist or was moved.",
          })}
        </p>
        <div className="flex items-center justify-center gap-3">
          <Button asChild variant="default">
            <Link to={homeHref}>
              {t("actions.go_home", { defaultValue: "Go to home" })}
            </Link>
          </Button>
          <Button
            variant="outline"
            onClick={() => window.history.back()}
          >
            {t("actions.go_back", { defaultValue: "Go back" })}
          </Button>
        </div>
      </div>
    </div>
  );
}
