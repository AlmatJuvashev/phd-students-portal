import React from 'react';
import { Link } from 'react-router-dom';
import { 
  GraduationCap, 
  ArrowRight, 
  BookOpen,
  Users,
  Calendar,
  FileText
} from 'lucide-react';
import { Button } from '@/components/ui/button';
import { useTranslation } from 'react-i18next';
import { APP_NAME } from '@/config';

export const Landing: React.FC = () => {
  const { t } = useTranslation("common");

  return (
    <div className="min-h-screen bg-background font-sans text-foreground selection:bg-primary/20 selection:text-primary flex flex-col">
      
      {/* Navbar */}
      <nav className="w-full bg-background/80 backdrop-blur-md border-b z-50 fixed top-0">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center gap-2">
              <div className="bg-primary p-1.5 rounded-lg">
                <GraduationCap className="h-5 w-5 text-primary-foreground" />
              </div>
              <span className="font-bold text-xl tracking-tight text-foreground">{APP_NAME}</span>
            </div>
            <div className="flex items-center gap-4">
              <Link to="/login">
                <Button variant="default" className="rounded-full group">
                  {t("landing.login", "Login")}
                  <ArrowRight className="h-4 w-4 ml-2 group-hover:translate-x-0.5 transition-transform" />
                </Button>
              </Link>
            </div>
          </div>
        </div>
      </nav>

      {/* Hero Section */}
      <div className="relative flex-1 flex flex-col justify-center items-center pt-24 overflow-hidden">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10 text-center">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-primary/10 border border-primary/20 text-primary text-sm font-medium mb-8 animate-in fade-in slide-in-from-bottom-4 duration-700">
            <span className="relative flex h-2 w-2">
              <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75"></span>
              <span className="relative inline-flex rounded-full h-2 w-2 bg-primary"></span>
            </span>
            {t("landing.badge", "KazNMU Doctoral Program")}
          </div>
          
          <h1 className="text-5xl md:text-7xl font-bold tracking-tight text-foreground mb-6 max-w-4xl mx-auto leading-[1.1] animate-in fade-in slide-in-from-bottom-6 duration-700 delay-100">
            {t("landing.title_prefix", "Streamlined Doctoral")} <br />
            <span className="text-transparent bg-clip-text bg-gradient-to-r from-primary to-blue-600">
              {t("landing.title_suffix", "Journey Management")}
            </span>
          </h1>
          
          <p className="text-lg md:text-xl text-muted-foreground mb-10 max-w-2xl mx-auto leading-relaxed animate-in fade-in slide-in-from-bottom-6 duration-700 delay-200">
            {t("landing.subtitle", "Navigate your PhD program with clarity. Track milestones, manage documents, and collaborate with advisors in one unified platform.")}
          </p>
          
          <div className="flex flex-col sm:flex-row gap-4 justify-center items-center animate-in fade-in slide-in-from-bottom-6 duration-700 delay-300">
            <Link to="/login">
              <Button size="lg" className="rounded-xl text-lg px-8 py-6 shadow-lg shadow-primary/20 hover:shadow-xl hover:shadow-primary/30 transition-all transform hover:-translate-y-0.5 w-full sm:w-auto">
                {t("landing.cta", "Start Your Journey")}
              </Button>
            </Link>
          </div>

          {/* Features Grid */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8 mt-20 max-w-5xl mx-auto text-left animate-in fade-in slide-in-from-bottom-8 duration-700 delay-500">
            <div className="p-6 rounded-2xl bg-card border shadow-sm hover:shadow-md transition-all">
              <div className="h-10 w-10 rounded-lg bg-blue-100 dark:bg-blue-900/30 flex items-center justify-center mb-4">
                <BookOpen className="h-5 w-5 text-blue-600 dark:text-blue-400" />
              </div>
              <h3 className="font-semibold text-lg mb-2">{t("landing.features.milestones.title", "Milestone Tracking")}</h3>
              <p className="text-muted-foreground">{t("landing.features.milestones.desc", "Visualize your progress through the doctoral program with clear milestones and deadlines.")}</p>
            </div>
            <div className="p-6 rounded-2xl bg-card border shadow-sm hover:shadow-md transition-all">
              <div className="h-10 w-10 rounded-lg bg-purple-100 dark:bg-purple-900/30 flex items-center justify-center mb-4">
                <FileText className="h-5 w-5 text-purple-600 dark:text-purple-400" />
              </div>
              <h3 className="font-semibold text-lg mb-2">{t("landing.features.documents.title", "Document Management")}</h3>
              <p className="text-muted-foreground">{t("landing.features.documents.desc", "Securely upload, organize, and share your research documents and administrative forms.")}</p>
            </div>
            <div className="p-6 rounded-2xl bg-card border shadow-sm hover:shadow-md transition-all">
              <div className="h-10 w-10 rounded-lg bg-teal-100 dark:bg-teal-900/30 flex items-center justify-center mb-4">
                <Users className="h-5 w-5 text-teal-600 dark:text-teal-400" />
              </div>
              <h3 className="font-semibold text-lg mb-2">{t("landing.features.collaboration.title", "Advisor Collaboration")}</h3>
              <p className="text-muted-foreground">{t("landing.features.collaboration.desc", "Seamlessly communicate with your advisors and committee members through integrated chat.")}</p>
            </div>
          </div>
        </div>

        {/* Abstract Background Decoration */}
        <div className="absolute top-0 left-1/2 -translate-x-1/2 w-full h-full -z-10 pointer-events-none overflow-hidden">
          <div className="absolute top-20 left-[10%] w-72 h-72 bg-primary/20 rounded-full blur-3xl mix-blend-multiply dark:mix-blend-screen filter opacity-70 animate-blob"></div>
          <div className="absolute top-20 right-[10%] w-72 h-72 bg-blue-400/20 rounded-full blur-3xl mix-blend-multiply dark:mix-blend-screen filter opacity-70 animate-blob animation-delay-2000"></div>
          <div className="absolute -bottom-8 left-[20%] w-72 h-72 bg-purple-400/20 rounded-full blur-3xl mix-blend-multiply dark:mix-blend-screen filter opacity-70 animate-blob animation-delay-4000"></div>
        </div>
      </div>
      
      {/* Simple Footer */}
      <div className="py-6 text-center text-muted-foreground text-sm border-t mt-auto">
        {t("landing.footer", { year: new Date().getFullYear(), defaultValue: "© KazNMU PhD Portal. All rights reserved." })}
      </div>
    </div>
  );
};

export default Landing;
