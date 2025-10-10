import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Link } from "react-router-dom";
import { Map, Users } from "lucide-react";
import { useTranslation } from "react-i18next";

export function HomePage() {
  const { t: T } = useTranslation("common");
  return (
    <div className="max-w-3xl mx-auto px-4 py-8 grid gap-4 md:grid-cols-2">
      <Card className="p-6 hover:shadow-md transition-shadow">
        <CardContent className="px-0 pb-0">
          <div className="flex items-start gap-3">
            <div className="rounded-xl bg-primary/10 p-3 text-primary">
              <Map className="h-6 w-6" />
            </div>
            <div className="flex-1">
              <div className="text-lg font-semibold mb-1">{T("nav.journey")}</div>
              <div className="text-sm text-muted-foreground mb-3">
                {T("home.journey_hint", { defaultValue: "Open your doctoral journey map." })}
              </div>
              <Link to="/journey">
                <Button>{T("home.open_journey", { defaultValue: "Open Journey" })}</Button>
              </Link>
            </div>
          </div>
        </CardContent>
      </Card>
      <Card className="p-6 hover:shadow-md transition-shadow">
        <CardContent className="px-0 pb-0">
          <div className="flex items-start gap-3">
            <div className="rounded-xl bg-primary/10 p-3 text-primary">
              <Users className="h-6 w-6" />
            </div>
            <div className="flex-1">
              <div className="text-lg font-semibold mb-1">{T("home.supervisors", { defaultValue: "Supervisors" })}</div>
              <div className="text-sm text-muted-foreground mb-3">
                {T("home.contacts_hint", { defaultValue: "Find and copy your supervisors’ contact details." })}
              </div>
              <Link to="/contacts">
                <Button variant="secondary">{T("home.open_contacts", { defaultValue: "Contacts" })}</Button>
              </Link>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

export default HomePage;

