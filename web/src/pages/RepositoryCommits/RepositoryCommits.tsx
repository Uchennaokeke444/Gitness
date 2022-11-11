import React from 'react'
import { Container, PageBody } from '@harness/uicore'
import { useGetRepositoryMetadata } from 'hooks/useGetRepositoryMetadata'
import { RepositoryCommitsContent } from './RepositoryCommitsContent/RepositoryCommitsContent'
import { RepositoryCommitsHeader } from './RepositoryCommitsHeader/RepositoryCommitsHeader'
import css from './RepositoryCommits.module.scss'

export default function RepositoryCommits() {
  const { repoMetadata, error, loading, commitRef } = useGetRepositoryMetadata()

  return (
    <Container className={css.main}>
      <PageBody loading={loading} error={error}>
        {repoMetadata ? (
          <>
            <RepositoryCommitsHeader repoMetadata={repoMetadata} />
            <RepositoryCommitsContent
              repoMetadata={repoMetadata}
              commitRef={commitRef || (repoMetadata.defaultBranch as string)}
            />
          </>
        ) : null}
      </PageBody>
    </Container>
  )
}