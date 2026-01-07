import React, { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { Loader2, QrCode, Keyboard } from 'lucide-react';
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription } from '@/components/ui/dialog';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { useToast } from '@/components/ui/use-toast';
import { Scanner } from '@yudiel/react-qr-scanner';
import { checkInStudent } from '../api';

interface CheckInModalProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

export const CheckInModal: React.FC<CheckInModalProps> = ({ open, onOpenChange }) => {
  const { t } = useTranslation('common');
  const { toast } = useToast();
  const qc = useQueryClient();
  const [sessionId, setSessionId] = useState('');
  const [activeTab, setActiveTab] = useState('scan');

  const checkInMutation = useMutation({
    mutationFn: (data: { session_id: string; code: string }) => checkInStudent(data),
    onSuccess: () => {
      toast({
        title: 'Checked In!',
        description: 'You have strictly marked your attendance.',
      });
      qc.invalidateQueries({ queryKey: ['student', 'dashboard'] });
      onOpenChange(false);
      setSessionId('');
    },
    onError: (err: any) => {
      toast({
        title: 'Check-in Failed',
        description: err.response?.data?.error || 'Could not verify code. Please try again.',
        variant: 'destructive',
      });
    },
  });

  const handleScan = (result: string) => {
    if (!result) return;
    try {
      // Expecting JSON: { "session_id": "...", "code": "..." }
      const data = JSON.parse(result);
      if (data.session_id) {
        checkInMutation.mutate({ session_id: data.session_id, code: data.code || 'QR' });
      }
    } catch (e) {
      // If not JSON, assume it's just the session ID
      checkInMutation.mutate({ session_id: result, code: 'QR' });
    }
  };

  const manualSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!sessionId) return;
    checkInMutation.mutate({ session_id: sessionId, code: 'MANUAL' });
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-md">
        <DialogHeader>
          <DialogTitle className="flex items-center gap-2">
            <QrCode className="h-5 w-5" /> Attendance Check-In
          </DialogTitle>
          <DialogDescription>
            Scan the QR code displayed by your instructor or enter the Session ID manually.
          </DialogDescription>
        </DialogHeader>

        <Tabs value={activeTab} onValueChange={setActiveTab} className="w-full">
          <TabsList className="grid w-full grid-cols-2">
            <TabsTrigger value="scan"><QrCode className="mr-2 h-4 w-4"/> Scan QR</TabsTrigger>
            <TabsTrigger value="manual"><Keyboard className="mr-2 h-4 w-4"/> Manual Entry</TabsTrigger>
          </TabsList>
          
          <TabsContent value="scan" className="space-y-4 py-4 min-h-[300px]">
            <div className="rounded-xl overflow-hidden border border-slate-200 aspect-square relative bg-black">
                <Scanner 
                    onScan={(result) => handleScan(result[0].rawValue)} 
                    allowMultiple={true}
                    scanDelay={2000}
                    components={{
                      finder: true,
                    }}
                />
            </div>
            <p className="text-xs text-center text-slate-500">
              Point your camera at the QR code on the teacher's screen.
            </p>
          </TabsContent>

          <TabsContent value="manual">
            <form onSubmit={manualSubmit} className="space-y-4 py-4">
              <div className="space-y-2">
                <Label htmlFor="session-id">Session ID</Label>
                <Input
                  id="session-id"
                  placeholder="e.g. sess-123..."
                  value={sessionId}
                  onChange={(e) => setSessionId(e.target.value)}
                  autoFocus
                />
              </div>
              
              <Button type="submit" className="w-full" disabled={!sessionId || checkInMutation.isPending}>
                {checkInMutation.isPending ? <Loader2 className="mr-2 h-4 w-4 animate-spin" /> : null}
                Check In
              </Button>
            </form>
          </TabsContent>
        </Tabs>
      </DialogContent>
    </Dialog>
  );
};
