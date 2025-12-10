import Link from "next/link"
import { MessageCircle, Settings, ArrowRight, Sparkles } from "lucide-react"
import { Button } from "@/components/ui/button"
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card"

export default function HomePage() {
  return (
    <div className="min-h-[100dvh] bg-gradient-to-b from-background via-background to-muted/30 flex items-center justify-center p-4 md:p-6">
      <div className="max-w-2xl w-full space-y-8">
        <div className="text-center space-y-3">
          <div className="inline-flex items-center gap-2 text-xs font-medium text-primary bg-primary/10 px-3 py-1 rounded-full mb-2">
            <Sparkles className="h-3 w-3" />
            Демо-версия
          </div>
          <h1 className="text-2xl md:text-3xl font-bold text-foreground text-balance">WebApp Checklist</h1>
          <p className="text-muted-foreground text-sm md:text-base">
            Модуль группового чата для медицинского университета
          </p>
        </div>

        <div className="grid gap-4 md:grid-cols-2">
          <Card className="group hover:shadow-lg transition-all duration-300 hover:-translate-y-1 border-primary/10 hover:border-primary/30">
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="p-2.5 rounded-xl bg-blue-100 group-hover:bg-blue-200 transition-colors">
                  <MessageCircle className="h-5 w-5 text-blue-600" />
                </div>
                <div>
                  <CardTitle className="text-lg">Чат</CardTitle>
                  <CardDescription>Сообщения и обсуждения</CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground mb-4 leading-relaxed">
                Общение с когортой, научными руководителями и администрацией программы.
              </p>
              <Button asChild className="w-full group-hover:bg-primary/90">
                <Link href="/chat">
                  Открыть чат
                  <ArrowRight className="h-4 w-4 ml-2 transition-transform group-hover:translate-x-1" />
                </Link>
              </Button>
            </CardContent>
          </Card>

          <Card className="group hover:shadow-lg transition-all duration-300 hover:-translate-y-1">
            <CardHeader>
              <div className="flex items-center gap-3">
                <div className="p-2.5 rounded-xl bg-slate-100 group-hover:bg-slate-200 transition-colors">
                  <Settings className="h-5 w-5 text-slate-600" />
                </div>
                <div>
                  <CardTitle className="text-lg">Администрирование</CardTitle>
                  <CardDescription>Управление группами</CardDescription>
                </div>
              </div>
            </CardHeader>
            <CardContent>
              <p className="text-sm text-muted-foreground mb-4 leading-relaxed">
                Создание групп, управление участниками и настройки чата.
              </p>
              <Button asChild variant="outline" className="w-full bg-transparent group-hover:bg-muted/50">
                <Link href="/admin/chat-rooms">
                  Управление
                  <ArrowRight className="h-4 w-4 ml-2 transition-transform group-hover:translate-x-1" />
                </Link>
              </Button>
            </CardContent>
          </Card>
        </div>

        <div className="text-center text-xs text-muted-foreground">
          <p>Демонстрационная версия с моковыми данными</p>
        </div>
      </div>
    </div>
  )
}
