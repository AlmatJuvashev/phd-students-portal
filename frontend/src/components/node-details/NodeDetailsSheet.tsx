// components/node-details/NodeDetailsSheet.tsx

import { NodeVM, t } from "@/lib/playbook";
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
} from "@/components/ui/sheet";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";

export function NodeDetailsSheet({
  node,
  onOpenChange,
}: {
  node: NodeVM | null;
  onOpenChange: (open: boolean) => void;
}) {
  return (
    <Sheet open={!!node} onOpenChange={onOpenChange}>
      <SheetContent side="right" className="w-[540px] sm:w-[600px]">
        {node && (
          <>
            <SheetHeader>
              <SheetTitle className="flex items-center gap-2">
                <span>{t(node.title, node.id)}</span>
                <Badge variant="secondary" className="capitalize">
                  {node.type}
                </Badge>
                <Badge className="capitalize">
                  {node.state.replace("_", " ")}
                </Badge>
              </SheetTitle>
            </SheetHeader>

            <div className="mt-6">
              <Tabs defaultValue="overview">
                <TabsList className="grid w-full grid-cols-4">
                  <TabsTrigger value="overview">Обзор</TabsTrigger>
                  <TabsTrigger value="todo">Что делать</TabsTrigger>
                  <TabsTrigger value="uploads">Загрузки</TabsTrigger>
                  <TabsTrigger value="activity">Активность</TabsTrigger>
                </TabsList>

                <TabsContent value="overview" className="space-y-4">
                  <Card className="p-4">
                    <div className="text-sm text-muted-foreground">
                      ID: {node.id}
                    </div>
                    {node.timer && (
                      <div className="text-sm">
                        Таймер: {node.timer.duration_days} дней
                      </div>
                    )}
                    <div className="text-sm">
                      Исполнители: {node.who_can_complete.join(", ")}
                    </div>
                  </Card>

                  <div className="flex gap-2">
                    <Button>Открыть шаг</Button>
                    <Button variant="secondary">Отметить выполненным</Button>
                  </div>
                </TabsContent>

                <TabsContent value="todo">
                  <Card className="p-4 text-sm">
                    {node.requirements?.notes && (
                      <p className="mb-2">{node.requirements.notes}</p>
                    )}
                    {(node.requirements?.fields?.length ?? 0) === 0 &&
                    (node.requirements?.uploads?.length ?? 0) === 0 ? (
                      <p>Для этого шага нет явных требований.</p>
                    ) : (
                      <div className="space-y-3">
                        {node.requirements?.fields?.length ? (
                          <div>
                            <div className="mb-1 font-medium">Поля формы</div>
                            <ul className="list-inside list-disc">
                              {node.requirements.fields.map((f) => (
                                <li key={f.key}>
                                  {f.key} {f.required ? "(обязательно)" : ""}
                                </li>
                              ))}
                            </ul>
                          </div>
                        ) : null}
                        {node.requirements?.validations?.length ? (
                          <div>
                            <div className="mb-1 font-medium">Валидации</div>
                            <ul className="list-inside list-disc">
                              {node.requirements.validations.map((v, i) => (
                                <li key={i}>
                                  {v.rule}
                                  {v.source ? ` @ ${v.source}` : ""}
                                </li>
                              ))}
                            </ul>
                          </div>
                        ) : null}
                      </div>
                    )}
                  </Card>
                </TabsContent>

                <TabsContent value="uploads">
                  <Card className="p-4 text-sm">
                    {node.requirements?.uploads?.length ? (
                      <ul className="list-inside list-disc">
                        {node.requirements.uploads.map((u) => (
                          <li key={u.key}>
                            {u.key} {u.required ? "(обязательно)" : ""}{" "}
                            {u.mime?.length ? `— ${u.mime.join(", ")}` : ""}
                          </li>
                        ))}
                      </ul>
                    ) : (
                      <p>Загрузки не требуются.</p>
                    )}
                  </Card>
                </TabsContent>

                <TabsContent value="activity">
                  <Card className="p-4 text-sm">
                    <p>
                      История действий появится здесь (сдачи, проверки,
                      статусы).
                    </p>
                  </Card>
                </TabsContent>
              </Tabs>
            </div>
          </>
        )}
      </SheetContent>
    </Sheet>
  );
}
