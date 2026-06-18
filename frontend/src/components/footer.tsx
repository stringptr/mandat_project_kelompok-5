interface FooterLink {
    label: string;
    href: string;
}

const FOOTER_LINKS: FooterLink[] = [
    { label: 'Tentang', href: '#' },
    { label: 'Kebijakan Privasi', href: '#' },
    { label: 'Kontak', href: '#' },
];

export function Footer(): JSX.Element {
    const appVersion: string = '2.4.0-Stable';

    return (
        <footer className="flex items-center justify-between px-8 py-4 bg-white border-t border-neutral-100">
            <div className="text-sm text-neutral-500 font-body">
                <span>© {new Date().getFullYear()} Dinas Kesehatan - SiGizi Komunitas. </span>
                <span className="text-neutral-400">Versi {appVersion}</span>
            </div>

            <nav className="flex items-center gap-6">
                {FOOTER_LINKS.map((link: FooterLink) => (
                    <a key={link.label} href={link.href} className="text-sm text-neutral-600 hover:text-primary transition-colors font-body">
                        {link.label}
                    </a>
                ))}
            </nav>
        </footer>
    );
}