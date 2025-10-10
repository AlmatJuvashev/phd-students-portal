import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Link } from "react-router-dom";
import { Map, Users } from "lucide-react";
import { useTranslation } from "react-i18next";

export function HomePage() {
  const { t: T } = useTranslation("common");
  return (
    <div className="max-w-4xl mx-auto px-4 py-8 space-y-6">
      {/* Hero Section */}
      <div className="text-center space-y-3 py-8">
        <h1 className="text-3xl sm:text-4xl font-bold bg-gradient-to-r from-primary via-primary to-primary/60 bg-clip-text text-transparent">
          PhD Student Portal
        </h1>
        <p className="text-base sm:text-lg text-muted-foreground max-w-2xl mx-auto">
          {T("home.welcome", { defaultValue: "Navigate your doctoral journey with clarity and confidence" })}
        </p>
      </div>

      {/* Action Cards */}
      <div className="grid gap-6 md:grid-cols-2">
        <Card className="bg-gradient-to-br from-card to-card/50 border-2 hover:border-primary/40 transition-all duration-300 hover:shadow-xl group">
          <CardContent className="p-6 sm:p-8">
            <div className="flex items-start gap-4">
              <div className="rounded-2xl bg-gradient-to-br from-primary to-primary/80 p-4 text-white shadow-lg group-hover:scale-110 transition-transform duration-300">
                <Map className="h-7 w-7" />
              </div>
              <div className="flex-1 space-y-3">
                <div className="text-xl font-bold">{T("nav.journey")}</div>
                <div className="text-sm text-muted-foreground leading-relaxed">
                  {T("home.journey_hint", { defaultValue: "Open your doctoral journey map and track your progress through each milestone." })}
                </div>
                <Link to="/journey">
                  <Button size="lg" className="mt-2 w-full sm:w-auto group-hover:shadow-lg transition-shadow">
                    {T("home.open_journey", { defaultValue: "Open Journey" })}
                    <Map className="ml-2 w-4 h-4" />
                  </Button>
                </Link>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card className="bg-gradient-to-br from-card to-card/50 border-2 hover:border-primary/40 transition-all duration-300 hover:shadow-xl group">
          <CardContent className="p-6 sm:p-8">
            <div className="flex items-start gap-4">
              <div className="rounded-2xl bg-gradient-to-br from-secondary to-secondary/80 p-4 text-secondary-foreground shadow-lg group-hover:scale-110 transition-transform duration-300">
                <Users className="h-7 w-7" />
              </div>
              <div className="flex-1 space-y-3">
                <div className="text-xl font-bold">{T("home.supervisors", { defaultValue: "Supervisors" })}</div>
                <div className="text-sm text-muted-foreground leading-relaxed">
                  {T("home.contacts_hint", { defaultValue: "Find and copy your supervisors' contact details quickly." })}
                </div>
                <Link to="/contacts">
                  <Button variant="secondary" size="lg" className="mt-2 w-full sm:w-auto group-hover:shadow-lg transition-shadow">
                    {T("home.open_contacts", { defaultValue: "View Contacts" })}
                    <Users className="ml-2 w-4 h-4" />
                  </Button>
                </Link>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

export default HomePage;
