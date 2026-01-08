
import React from 'react';
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

export const HRPage = () => {
    return (
        <div className="p-8">
            <h1 className="text-2xl font-bold mb-6">HR Management</h1>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <Card>
                    <CardHeader>
                        <CardTitle>Staff Directory</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <p className="text-muted-foreground">Manage academic and administrative staff records.</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader>
                        <CardTitle>Recruitment</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <p className="text-muted-foreground">Active job postings and candidate tracking.</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader>
                        <CardTitle>Payroll & Benefits</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <p className="text-muted-foreground">Manage salary periods and benefit packages.</p>
                    </CardContent>
                </Card>
            </div>
        </div>
    );
};

export default HRPage;
