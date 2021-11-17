import { Component, h, RenderableProps } from 'react';
import * as style from './badge.css';

interface BadgeProps {
    type: string;
    className?: string;
}

class Badge extends Component<BadgeProps> {
    public render(
        { type, className, children }: RenderableProps<BadgeProps>,
    ) {
        return (
            <span className={[style.container, style[type], className].join(' ')}>
                {children}
            </span>
        );
    }
}

export { Badge };
