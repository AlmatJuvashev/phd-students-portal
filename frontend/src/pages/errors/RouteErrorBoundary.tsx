import { isRouteErrorResponse, useLocation, useRouteError } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/button";
import { NotFoundIllustration } from "./NotFoundIllustration";
import { Link } from "react-router-dom";

export default function RouteErrorBoundary() {
  const { t } = useTranslation("common");
  const error = useRouteError();
  const loc = useLocation();
  const isAdmin = loc.pathname.startsWith("/admin");
  const homeHref = isAdmin ? "/admin" : "/";

  if (isRouteErrorResponse(error) && error.status === 404) {
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
            <Button variant="outline" onClick={() => window.history.back()}>
              {t("actions.go_back", { defaultValue: "Go back" })}
            </Button>
          </div>
        </div>
      </div>
    );
  }

  const title = t("errors.unexpected_title", { defaultValue: "Something went wrong" });
  const desc = t("errors.unexpected_desc", {
    defaultValue: "An unexpected error occurred. Please try again.",
  });

  return (
    <div className="min-h-[60vh] flex items-center justify-center px-4">
      <div className="max-w-2xl text-center space-y-5">
        <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-muted">
          <span className="text-2xl font-bold">⚠️</span>
        </div>
        <h1 className="text-2xl sm:text-3xl font-bold">{title}</h1>
        <p className="text-muted-foreground">{desc}</p>
        {process.env.NODE_ENV !== "production" && (
          <pre className="text-xs text-left whitespace-pre-wrap bg-muted p-3 rounded-md overflow-auto max-h-48">
            {String((error as any)?.statusText || (error as any)?.message || error)}
          </pre>
        )}
        <div className="flex items-center justify-center gap-3">
          <Button onClick={() => window.location.reload()}>
            {t("actions.reload", { defaultValue: "Reload" })}
          </Button>
          <Button asChild variant="outline">
            <Link to={homeHref}>
              {t("actions.go_home", { defaultValue: "Go to home" })}
            </Link>
          </Button>
        </div>
      </div>
    </div>
  );
}
