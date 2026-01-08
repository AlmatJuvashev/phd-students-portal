
import React from 'react';
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Building, MapPin, Wrench } from 'lucide-react';

export const FacilitiesPage = () => {
    return (
        <div className="p-8">
            <h1 className="text-2xl font-bold mb-6">Facilities Management</h1>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
                <Card>
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2"><Building size={20} /> Campus Map</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <p className="text-muted-foreground">Interactive map of all campus buildings and rooms.</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2"><MapPin size={20} /> Room Utilization</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <p className="text-muted-foreground">Heatmaps and analytics for room bookings.</p>
                    </CardContent>
                </Card>
                <Card>
                    <CardHeader>
                        <CardTitle className="flex items-center gap-2"><Wrench size={20} /> Maintenance</CardTitle>
                    </CardHeader>
                    <CardContent>
                        <p className="text-muted-foreground">Track repair requests and scheduled maintenance.</p>
                    </CardContent>
                </Card>
            </div>
        </div>
    );
};

export default FacilitiesPage;
