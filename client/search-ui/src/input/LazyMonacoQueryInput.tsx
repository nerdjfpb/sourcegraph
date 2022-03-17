import classNames from 'classnames'
import React, { Suspense } from 'react'

import { SettingsExperimentalFeatures } from '@sourcegraph/shared/src/schema/settings.schema'
import { lazyComponent } from '@sourcegraph/shared/src/util/lazyComponent'

import styles from './LazyMonacoQueryInput.module.scss'
import { MonacoQueryInputProps } from './MonacoQueryInput'

/**
 * Minimal interface for external interaction with the editor.
 */
export interface IEditor {
    focus(): void
}

const MonacoQueryInput = lazyComponent(() => import('./MonacoQueryInput'), 'MonacoQueryInput')
const CodemirrorQueryInput = lazyComponent(() => import('./CodemirrorQueryInput'), 'CodemirrorMonacoFacade')

/**
 * A plain query input displayed during lazy-loading of the MonacoQueryInput.
 * It has no suggestions, but still allows to type in and submit queries.
 */
export const PlainQueryInput: React.FunctionComponent<
    Pick<MonacoQueryInputProps, 'queryState' | 'autoFocus' | 'onChange' | 'className'>
> = ({ queryState, autoFocus, onChange, className }) => {
    const onInputChange = React.useCallback(
        (event: React.ChangeEvent<HTMLInputElement>) => {
            onChange({ query: event.target.value })
        },
        [onChange]
    )
    return (
        <input
            type="text"
            autoFocus={autoFocus}
            className={classNames('form-control text-code', styles.lazyMonacoQueryInputIntermediateInput, className)}
            value={queryState.query}
            onChange={onInputChange}
            spellCheck={false}
        />
    )
}

export interface LazyMonacoQueryInputProps extends MonacoQueryInputProps {
    /**
     * Determines which editor implementation to use.
     */
    editorComponent: SettingsExperimentalFeatures['editor']
}

/**
 * A lazily-loaded {@link MonacoQueryInput}, displaying a read-only query field as a fallback during loading.
 */
export const LazyMonacoQueryInput: React.FunctionComponent<LazyMonacoQueryInputProps> = ({
    editorComponent,
    ...props
}) => {
    const QueryInput = editorComponent === 'codemirror6' ? CodemirrorQueryInput : MonacoQueryInput

    return (
        <Suspense fallback={<PlainQueryInput {...props} />}>
            <QueryInput {...props} />
        </Suspense>
    )
}
