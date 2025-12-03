import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Link } from "react-router-dom";
import { Map, Users, MessageCircle, Calendar } from "lucide-react";
import { useTranslation } from "react-i18next";

export function HomePage() {
  const { t: T } = useTranslation("common");
  return (
    <div className="max-w-6xl mx-auto px-4 py-8 space-y-8">
      {/* Hero Section */}
      <div className="text-center space-y-4 py-8">
        <h1 className="text-3xl sm:text-5xl font-bold bg-gradient-to-r from-primary via-primary to-primary/60 bg-clip-text text-transparent">
          {T("home.hero_title", { defaultValue: "PhD Student Portal" })}
        </h1>
        <p className="text-base sm:text-xl text-muted-foreground max-w-2xl mx-auto">
          {T("home.hero_subtitle", {
            defaultValue:
              "Navigate your doctoral journey with clarity and confidence",
          })}
        </p>
      </div>

      {/* Action Cards */}
      <div className="grid gap-6 md:grid-cols-2 lg:grid-cols-4">
        <Card className="bg-gradient-to-br from-card to-card/50 border-2 hover:border-primary/40 transition-all duration-300 hover:shadow-xl group h-full">
          <CardContent className="p-6 flex flex-col items-center text-center gap-4 h-full">
            <div className="rounded-2xl bg-gradient-to-br from-primary to-primary/80 p-4 text-white shadow-lg group-hover:scale-110 transition-transform duration-300">
              <Map className="h-8 w-8" />
            </div>
            <div className="space-y-2 flex-1">
              <div className="text-xl font-bold">{T("nav.journey")}</div>
              <p className="text-sm text-muted-foreground">
                {T("home.journey_hint", {
                  defaultValue: "Track your progress through milestones.",
                })}
              </p>
            </div>
            <Link to="/journey" className="w-full">
              <Button className="w-full group-hover:shadow-lg transition-shadow">
                {T("home.open_journey", { defaultValue: "Open Journey" })}
              </Button>
            </Link>
          </CardContent>
        </Card>

        <Card className="bg-gradient-to-br from-card to-card/50 border-2 hover:border-primary/40 transition-all duration-300 hover:shadow-xl group h-full">
          <CardContent className="p-6 flex flex-col items-center text-center gap-4 h-full">
            <div className="rounded-2xl bg-gradient-to-br from-secondary to-secondary/80 p-4 text-secondary-foreground shadow-lg group-hover:scale-110 transition-transform duration-300">
              <Users className="h-8 w-8" />
            </div>
            <div className="space-y-2 flex-1">
              <div className="text-xl font-bold">
                {T("home.supervisors", { defaultValue: "Supervisors" })}
              </div>
              <p className="text-sm text-muted-foreground">
                {T("home.contacts_hint", {
                  defaultValue: "Find supervisors' contact details.",
                })}
              </p>
            </div>
            <Link to="/contacts" className="w-full">
              <Button
                variant="secondary"
                className="w-full group-hover:shadow-lg transition-shadow"
              >
                {T("home.open_contacts", { defaultValue: "View Contacts" })}
              </Button>
            </Link>
          </CardContent>
        </Card>

        <Card className="bg-gradient-to-br from-card to-card/50 border-2 hover:border-primary/40 transition-all duration-300 hover:shadow-xl group h-full">
          <CardContent className="p-6 flex flex-col items-center text-center gap-4 h-full">
            <div className="rounded-2xl bg-gradient-to-br from-blue-500 to-blue-600 p-4 text-white shadow-lg group-hover:scale-110 transition-transform duration-300">
              <MessageCircle className="h-8 w-8" />
            </div>
            <div className="space-y-2 flex-1">
              <div className="text-xl font-bold">
                {T("nav.chat", { defaultValue: "Messages" })}
              </div>
              <p className="text-sm text-muted-foreground">
                Chat with your advisors and colleagues.
              </p>
            </div>
            <Link to="/chat" className="w-full">
              <Button
                variant="outline"
                className="w-full group-hover:shadow-lg transition-shadow"
              >
                Open Chat
              </Button>
            </Link>
          </CardContent>
        </Card>

        <Card className="bg-gradient-to-br from-card to-card/50 border-2 hover:border-primary/40 transition-all duration-300 hover:shadow-xl group h-full">
          <CardContent className="p-6 flex flex-col items-center text-center gap-4 h-full">
            <div className="rounded-2xl bg-gradient-to-br from-purple-500 to-purple-600 p-4 text-white shadow-lg group-hover:scale-110 transition-transform duration-300">
              <Calendar className="h-8 w-8" />
            </div>
            <div className="space-y-2 flex-1">
              <div className="text-xl font-bold">
                {T("nav.calendar", { defaultValue: "Calendar" })}
              </div>
              <p className="text-sm text-muted-foreground">
                View upcoming events and deadlines.
              </p>
            </div>
            <Link to="/calendar" className="w-full">
              <Button
                variant="outline"
                className="w-full group-hover:shadow-lg transition-shadow"
              >
                View Calendar
              </Button>
            </Link>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

export default HomePage;
