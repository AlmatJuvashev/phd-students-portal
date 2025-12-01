import React from "react";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { ProgramsList } from "./ProgramsList";
import { SpecialtiesList } from "./SpecialtiesList";
import { CohortsList } from "./CohortsList";
import { DepartmentsList } from "./DepartmentsList";
import { useTranslation } from "react-i18next";

export const DictionariesPage = () => {
  const { t } = useTranslation();

  return (
    <div className="container mx-auto py-8">
      <h1 className="text-2xl font-bold mb-6">Dictionaries</h1>
      <Tabs defaultValue="programs" className="w-full">
        <TabsList className="mb-4">
          <TabsTrigger value="programs">Programs</TabsTrigger>
          <TabsTrigger value="specialties">Specialties</TabsTrigger>
          <TabsTrigger value="cohorts">Cohorts</TabsTrigger>
          <TabsTrigger value="departments">Departments</TabsTrigger>
        </TabsList>
        <TabsContent value="programs">
          <ProgramsList />
        </TabsContent>
        <TabsContent value="specialties">
          <SpecialtiesList />
        </TabsContent>
        <TabsContent value="cohorts">
          <CohortsList />
        </TabsContent>
        <TabsContent value="departments">
          <DepartmentsList />
        </TabsContent>
      </Tabs>
    </div>
  );
};

export default DictionariesPage;
