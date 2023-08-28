import React, { useMemo } from 'react'
import { orderBy } from 'lodash-es'
import { Container, TableV2 as Table, Text, Layout, Button, ButtonVariation } from '@harnessio/uicore'
import { Icon } from '@harnessio/icons'
import { Color } from '@harnessio/design-system'
import type { CellProps, Column } from 'react-table'
import { useStrings } from 'framework/strings'
import type { GitInfoProps } from 'utils/GitUtils'
import css from './RepositorySettings.module.scss'

interface Hook {
  url: string
}
interface SettingsContentProps extends Pick<GitInfoProps, 'repoMetadata'> {
  hooks: Hook[]
}

export function SettingsContent({ hooks }: SettingsContentProps) {
  const { getString } = useStrings()
  const columns: Column<Hook>[] = useMemo(
    () => [
      {
        id: 'url',
        width: '85%',
        Cell: ({ row }: CellProps<Hook>) => {
          return (
            <Layout.Horizontal spacing={'medium'}>
              <Icon name="deployment-success-legacy" />
              <Text intent={'primary'} lineClamp={1}>
                {row.original.url}
              </Text>
              <Text color={Color.BLACK}>({getString('webhookListingContent')})</Text>
            </Layout.Horizontal>
          )
        }
      },
      {
        id: 'actions',
        width: '15%',
        Cell: () => {
          return (
            <Layout.Horizontal flex>
              <Button variation={ButtonVariation.SECONDARY} intent="primary">
                {getString('edit')}
              </Button>
              <Button variation={ButtonVariation.SECONDARY} intent="danger">
                {getString('delete')}
              </Button>
            </Layout.Horizontal>
          )
        }
      }
    ],
    [getString]
  )

  return (
    <Container className={css.contentContainer}>
      <Table<Hook> hideHeaders columns={columns} data={orderBy(hooks)} />
    </Container>
  )
}