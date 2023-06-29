[Setup]
AppName=aliae
AppVersion=<VERSION>
DefaultDirName={autopf}\aliae
DefaultGroupName=aliae
AppPublisher=Jan De Dobbeleer
AppPublisherURL=https://aliae.dev
AppSupportURL=https://github.com/JanDeDobbeleer/aliae/issues
LicenseFile="bin\LICENSE.txt"
OutputBaseFilename=install
PrivilegesRequired=lowest
PrivilegesRequiredOverridesAllowed=dialog
ChangesEnvironment=yes
SignTool=signtool
SignedUninstaller=yes
CloseApplications=no

[Files]
Source: "bin\aliae.exe"; DestDir: "{app}\bin"; Flags: sign

[Registry]
Root: "HKA"; Subkey: "{code:GetEnvironmentKey}"; ValueType: expandsz; ValueName: "Path"; ValueData: "{olddata};{app}\bin"; Check: NeedsAddPathHKA(ExpandConstant('{app}\bin'))
Root: "HKA"; Subkey: "{code:GetEnvironmentKey}"; ValueType: string; ValueName: "ALIAE_INSTALLER"; ValueData: {param:installer|manual}; Flags: preservestringtype

[Code]
function GetEnvironmentKey(Param: string): string;
begin
  if IsAdminInstallMode then
    Result := 'System\CurrentControlSet\Control\Session Manager\Environment'
  else
    Result := 'Environment';
end;

function NeedsAddPathHKA(Param: string): boolean;
var
    OrigPath: string;
begin
    if not RegQueryStringValue(HKA, GetEnvironmentKey(''), 'Path', OrigPath)
    then begin
        Result := True;
        exit;
    end;
    // look for the path with leading and trailing semicolon
    // Pos() returns 0 if not found
    Result := Pos(';' + Param + ';', ';' + OrigPath + ';') = 0;
end;
